package p2p

import (
	"context"
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/crypto"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	"go.uber.org/zap"
	"io"
	"math"
	"math/rand"
	"net"
	"reflect"
	"sync"
	"time"
)

type HostConfig struct {
	PeerID   uint64
	BindAddr string
	Network  *config.Network

	ListenConfig *net.ListenConfig
}

type Node struct {
	Config     HostConfig
	Blockchain *cryptonote.BlockChain

	dialer *net.Dialer
	logger *zap.SugaredLogger
	wg     *sync.WaitGroup
	ps     *peerStore

	context context.Context

	listener *net.TCPListener
}

// NewNode creates instance of the node
func NewNode(core *cryptonote.BlockChain, cfg HostConfig, logger *zap.Logger) Node {
	var wg sync.WaitGroup

	h := Node{
		Config:     cfg,
		Blockchain: core,
		logger:     logger.Sugar(),
	}

	h.defaults()
	h.ps = NewPeerStore()
	h.wg = &wg

	return h
}

// Run the node server
// Listen for new connections on the defined port.
// Send handshake requests to the seed nodes.
func (n *Node) Run(ctx context.Context) error {
	// listener, err := n.Config.ListenConfig.Listen(ctx, "tcp", n.Config.BindAddr)
	addr, err := net.ResolveTCPAddr("tcp", n.Config.BindAddr)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	n.listener = listener
	n.logger.Debugf("listening on %s", listener.Addr())

	n.context = ctx

	n.wg.Add(1)
	go n.runListener()

	// TODO: Disabled only for check the listener
	//for _, seedAddr := range n.Config.Network.SeedNodes {
	//	go n.syncWithAddr(seedAddr)
	//}

	n.wg.Wait()
	return nil
}

// runListener listens for a new connections from external peers.
func (n *Node) runListener() {
	for {
		select {
		case <-n.context.Done():
			if err := n.listener.Close(); err != nil {
				n.logger.Errorf("failed to close listener: %s", err)
			}

			n.wg.Done()
			return
		default:
			_ = n.listener.SetDeadline(time.Now().Add(time.Second * 5))

			conn, err := n.listener.AcceptTCP()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					time.Sleep(time.Second)
					continue
				}
				n.logger.Errorf("failed to accept connection: %s", err)
			}

			go n.handleIncomingConnection(conn)
		}
	}
}

func (n *Node) handleIncomingConnection(conn *net.TCPConn) {
	// TODO: Enabling handling incoming connections
	//return

	peer := NewPeerFromIncomingConnection(n, conn)

	//// TODO: Add peer to peerstore. Make sure it is not exists.
	//
	n.wg.Add(1)
	defer n.wg.Done()

	n.connectionHandler(peer)

	if err := n.ps.toGrey(peer); err != nil {
		n.logger.Warnf("peer remove failed: %s", err)
	}

	n.logger.Debugf("[%16x] sync closed", peer.ID)
}

func (n *Node) connectionHandler(p *Peer) {
	for {
		// Peer state changes asynchronously after handling some commands.
		// Here we are taking care of handle different peer statuses.
		switch p.state {
		// Our node must be synchronized with the peer
		case PeerStateSyncRequired:
			p.state = PeerStateSynchronizing

			if err := n.NotifyRequestChain(p); err != nil {
				n.logger.Errorf("failed to write request chain: %s", err)
			}

		// Peer shutdown.
		// Stop listening for commands.
		case PeerStateShutdown:
			n.logger.Infof("[%d] shutting down...", p.ID)
			return
		}

		// Wait 1 second before read next command from the connection
		// Stop listening for the commands by return from the function
		select {
		case <-time.After(time.Second):
		case <-n.context.Done():
			return
		}

		// Read a next command from the connection
		cmd, err := p.protocol.read()

		// No command received continue to next loop cycle
		if err == io.EOF {
			continue
		}

		// On any error we move the peer to the grey list
		if err != nil {
			n.logger.Errorf("error on read command: %s", err)
			_ = n.ps.toGrey(p)
			break
		}

		// There is special command type as "notification" we handle them with a separate method.
		if cmd.IsNotify {
			if err := n.handleNotification(p, cmd); err != nil {
				n.logger.Errorf("failed to handle notification %d: %s", cmd.Command, err)
			}

			continue
		}

		// Call method for handling the notification
		if err := n.handleCommand(p, cmd); err != nil {
			n.logger.Errorf("failed handle command (%d): %s", cmd.Command, err)
		}
	}
}

// handleNotification
//
// Receive notification from remote peer and handle it depend on the notification code.
func (n *Node) handleNotification(p *Peer, cmd *LevinCommand) error {
	p.logger.Debugf("handeling notification: %d", cmd.Command)

	// TODO: This part of the code can be used for debug
	//cwd, err := os.Getwd()
	//if err != nil {
	//	panic(err)
	//}
	//err = ioutil.WriteFile(
	//	fmt.Sprintf("%s/%d.dat", cwd, cmd.Command),
	//	cmd.Payload,
	//	0644,
	//)
	//if err != nil {
	//	panic(err)
	//}

	nt, err := parseNotification(cmd)
	if err != nil {
		return err
	}

	switch nt.(type) {
	case NotificationTxPool: // 2008
		notification := nt.(NotificationTxPool)

		p.logger.Debugf("notification tx pool, size: %d", len(notification.Transactions))
	case NotificationRequestChain:
		notification := nt.(NotificationRequestChain)

		// Verify more than 0 requested blocks
		if len(notification.Blocks) == 0 {
			p.Shutdown()
			return errors.New(fmt.Sprintf("[%s] request chain with 0 blocks", p))
		}

		genesisBlock, err := n.Blockchain.GenesisBlock()
		if err != nil {
			return fmt.Errorf("[%s] unexpected error: %w", p, err)
		}

		// Make sure genesis blocks belongs to same network
		if *genesisBlock.Hash() != notification.Blocks[len(notification.Blocks)-1] {
			p.Shutdown()
			return errors.New(fmt.Sprintf("[%s] request chain genesis block not match", p))
		}

		// Get max index of this node blockchain
		//topIndex, err := p.node.Blockchain.TopIndex()
		//if err != nil {
		//	return fmt.Errorf("[%s] unexpected error: %w", p, err)
		//}

		// TODO: Build response chain entry and response to requested peer
		//responseChainEntry := &NotificationResponseChainEntry{
		//	TotalHeight: topIndex + 1,
		//}

		n.logger.Debugf("[%s] request chain %d blocks.", p, len(notification.Blocks))
	case NotificationResponseChainEntry: // 2007
		notification := nt.(NotificationResponseChainEntry)

		p.logger.Debugf(
			"notification response chain entry, start: %d, blocks: %d, total: %d",
			notification.StartHeight, len(notification.BlocksHashes), notification.TotalHeight,
		)

		if len(notification.BlocksHashes) == 0 {
			p.Shutdown()
			// TODO: Create new error instance
			return errors.New(fmt.Sprintf("[%s] received empty blocks in response chain enrty", p))
		}

		firstHash := notification.BlocksHashes[0]
		hasFirstBlock := n.Blockchain.HaveBlock(&firstHash)

		if !hasFirstBlock {
			p.Shutdown()
			// TODO: Create new error instance
			return errors.New(fmt.Sprintf("[%s] hash %s missing in our blockchain", p, firstHash.String()))
		}

		p.remoteHeight = notification.TotalHeight
		p.lastResponseHeight = notification.StartHeight + uint32(len(notification.BlocksHashes)-1)

		if p.lastResponseHeight > p.remoteHeight {
			// TODO: Create new error instance
			p.Shutdown()
			return errors.New(
				fmt.Sprintf(
					"[%s] sent wrong response chain entry, with TotalHeight = %d, StartHeight = %d, blocks = %d", p,
					notification.StartHeight,
					notification.TotalHeight,
					len(notification.BlocksHashes),
				),
			)
		}

		allBlockKnown := true
		for _, bh := range notification.BlocksHashes {
			hasBlock := n.Blockchain.HaveBlock(&bh)

			if allBlockKnown && hasBlock {
				continue
			}

			allBlockKnown = false
			p.neededBlocks = append(p.neededBlocks, bh)
		}

		return p.requestMissingBlocks(n, false)
	case NotificationResponseGetObjects: // 2004
		notification := nt.(NotificationResponseGetObjects)

		return n.HandleResponseGetObjects(p, notification)
	default:
		n.logger.Errorf("can't handle notification type: %s", reflect.TypeOf(nt))
	}

	return nil
}

func (n *Node) handleCommand(p *Peer, cmd *LevinCommand) error {
	c, err := parseCommand(cmd)
	if err != nil {
		return err
	}

	switch c.(type) {
	case HandshakeRequest:
		if err := n.HandleHandshake(p, c.(HandshakeRequest)); err != nil {
			return err
		}

		rsp := NewHandshakeResponse(n.Blockchain, n.ps.toPeerEntries())
		if err := p.protocol.Reply(cmd.Command, rsp, 1); err != nil {
			return err
		}
	case TimedSyncRequest:
		command := c.(TimedSyncRequest)
		if err := n.processSyncData(p, command.PayloadData, false); err != nil {
			return err
		}

		res, err := newTimedSyncResponse(n)
		if err != nil {
			return err
		}

		if err := p.protocol.Reply(cmd.Command, *res, 1); err != nil {
			return err
		}

		n.logger.Infof("[%s] sync request %d", p, command.PayloadData.CurrentHeight)
	default:
		n.logger.Errorf("received unknown commands type: %s", reflect.TypeOf(c))
	}

	return nil
}

// processSyncData processing remote sync data
//
// This method is safe for concurrent calls.
func (n *Node) processSyncData(p *Peer, syncData SyncData, isInitial bool) error {
	p.Lock()
	defer p.Unlock()

	// Ignore all not initial requests
	if p.state == PeerStateBeforeHandshake && !isInitial {
		return nil
	}

	if p.state == PeerStateSynchronizing {
	} else if n.Blockchain.HaveBlock(&syncData.TopBlockHash) {
		if isInitial {
			n.onSynchronized()
			p.state = PeerStatePoolSyncRequired
		} else {
			p.state = PeerStateNormal
		}
	} else {
		height := n.Blockchain.Height()

		diff := int64(syncData.CurrentHeight) - int64(height)
		if diff < 0 && uint32(math.Abs(float64(diff))) > n.Blockchain.Network.MinedMoneyUnlockWindow() {
			if n.Blockchain.Checkpoints.IsInCheckpointZone(syncData.CurrentHeight) {

				p.logger.Debugf(
					"Sync data return a new top block candidate: %d -> %d\nYour node is %d blocks ahead.\n"+
						"The block candidate is too deep behind and in checkpoint zone, dropping connection",
					height,
					syncData.CurrentHeight,
					uint32(math.Abs(float64(diff))),
				)

				n.addHostFail(p.address)
				p.Shutdown()

				return ErrSyncDataTooDeepBehind
			}
		}

		p.logger.Infof(
			"Sync data returned a new top block candidate: %d -> %d\n"+
				"Your node is %d blocks behind/ahead. Synchronization started.",
			height,
			syncData.CurrentHeight,
			uint32(math.Abs(float64(diff))),
		)

		p.state = PeerStateSyncRequired
	}

	n.updateObservedHeight(p, syncData.CurrentHeight)
	p.remoteHeight = syncData.CurrentHeight

	// TODO: Implement notification
	// if (is_initial) {
	// 	 m_peersCount++;
	//	 m_observerManager.notify(&ICryptoNoteProtocolObserver::peerCountUpdated, m_peersCount.load());
	// }

	return nil
}

func (n *Node) processNewObjects(blocks []*cryptonote.Block, transactions map[crypto.Hash][][]byte) error {
	for i, block := range blocks {
		if _, ok := transactions[*block.Hash()]; !ok {
			return errors.New(fmt.Sprintf("transactions for block at index %d not found", i))
		}

		if err := n.Blockchain.AddBlock(block, transactions[*block.Hash()]); err != nil {
			return err
			// TODO: Process proper error
			//
			//if (addResult == error::AddBlockErrorCondition::BLOCK_VALIDATION_FAILED ||
			//	addResult == error::AddBlockErrorCondition::TRANSACTION_VALIDATION_FAILED ||
			//	addResult == error::AddBlockErrorCondition::DESERIALIZATION_FAILED) {
			//	logger(Logging::DEBUGGING) << context << "Block verification failed, dropping connection: " << addResult.message();
			//	m_p2p->drop_connection(context, true);
			//	return 1;
			//} else if (addResult == error::AddBlockErrorCondition::BLOCK_REJECTED) {
			//	logger(Logging::DEBUGGING) << context << "Block received at sync phase was marked as orphaned, dropping connection: " << addResult.message();
			//	m_p2p->drop_connection(context, true);
			//	return 1;
			//} else if (addResult == error::AddBlockErrorCode::ALREADY_EXISTS) {
			//	logger(Logging::DEBUGGING) << context << "Block already exists, switching to idle state: " << addResult.message();
			//	context.m_state = CryptoNoteConnectionContext::state_idle;
			//	context.m_needed_objects.clear();
			//	context.m_requested_objects.clear();
			//	return 1;
			//}
		}
	}

	return nil
}

// TODO: Implement CryptoNoteProtocolHandler::updateObservedHeight
func (n *Node) updateObservedHeight(p *Peer, height uint32) {
}

// TODO: Implement NodeServer::add_host_fail and rename this method
func (n *Node) addHostFail(address NetworkAddress) {
}

// TODO: Implement on_connection_synchronized
func (n *Node) onSynchronized() {
	// CryptoNoteProtocolHandler::on_connection_synchronized
	// LINE: 916
}

func (n *Node) defaults() {
	if n.Config.PeerID == 0 {
		n.Config.PeerID = rand.Uint64()
	}

	if n.Config.ListenConfig == nil {
		n.Config.ListenConfig = &net.ListenConfig{}
	}

	if n.dialer == nil {
		n.dialer = &net.Dialer{
			//LocalAddr: n.Config.BindAddr,
			Timeout: time.Second,
		}
	}
}

func (n *Node) syncWithAddr(addr string) {
	ctx, cancel := context.WithCancel(n.context)
	defer cancel()

	peer, err := NewPeerFromTCPAddress(ctx, n, addr)
	if err != nil {
		// n.logger.Errorf("failed to dial to peer: %s", err)
		cancel()
		return
	}

	//handshake, err := peer.handshake(n)
	_, err = peer.handshake(n)
	if err != nil {
		n.logger.Errorf("failed handshake: %s", err)
		cancel()
		return
	}

	n.logger.Debugf("[%s] handshake established", peer)

	if err := n.ps.toWhite(peer); err != nil {
		n.logger.Error("failed to add peer to the store")
		cancel()
		return
	}

	n.wg.Add(1)
	defer n.wg.Done()

	//for _, pe := range handshake.Peers {
	//	go n.syncWithAddr(c, pe.Address.String())
	//}

	n.connectionHandler(peer)

	if err := n.ps.toGrey(peer); err != nil {
		n.logger.Warnf("peer remove failed: %s", err)
	}

	n.logger.Debugf("[%16x] sync closed", peer.ID)
}
