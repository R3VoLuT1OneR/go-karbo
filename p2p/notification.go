package p2p

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/crypto"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	"github.com/r3volut1oner/go-karbo/encoding/binary"
	"reflect"
)

const (
	NotificationBoolBase             = 2000                      // 2000
	NotificationNewBlockID           = NotificationBoolBase + 1  // 2001
	NotificationNewTransactionsID    = NotificationBoolBase + 2  // 2002
	NotificationRequestGetObjectsID  = NotificationBoolBase + 3  // 2003
	NotificationResponseGetObjectsID = NotificationBoolBase + 4  // 2004
	NotificationRequestChainID       = NotificationBoolBase + 6  // 2006
	NotificationResponseChainEntryID = NotificationBoolBase + 7  // 2007
	NotificationTxPoolID             = NotificationBoolBase + 8  // 2008
	NotificationNewLiteBlockID       = NotificationBoolBase + 9  // 2009
	NotificationMissingTxsID         = NotificationBoolBase + 10 // 2010
)

type RawBlock struct {
	Block        []byte   `binary:"block"`
	Transactions [][]byte `binary:"txs,array,omitempty"`
}

type NotificationNewBlock struct {
	Block                   RawBlock `binary:"b"`
	CurrentBlockchainHeight uint32   `binary:"current_blockchain_height"`
	Hop                     uint32   `binary:"hop"`
}

type NotificationNewTransactions struct {
	Stem         bool          `binary:"stem"`
	Transactions []crypto.Hash `binary:"txs,binary"`
}

type NotificationRequestGetObjects struct {
	Transactions []crypto.Hash `binary:"txs,binary"`
	Blocks       []crypto.Hash `binary:"blocks,binary"`
}

type NotificationTxPool struct {
	Transactions []crypto.Hash `binary:"txs"`
}

type NotificationNewLiteBlock struct {
	CurrentBlockchainHeight uint32 `binary:"current_blockchain_height"`
	Hop                     uint32 `binary:"hop"`
	Block                   []byte `binary:"block"`
}

// NotificationResponseChainEntry = 2007
type NotificationResponseChainEntry struct {
	StartHeight  uint32        `binary:"start_height"`
	TotalHeight  uint32        `binary:"total_height"`
	BlocksHashes []crypto.Hash `binary:"m_block_ids,binary"`
}

var mapNotificationID = map[uint32]interface{}{
	NotificationNewBlockID:           NotificationNewBlock{},
	NotificationNewTransactionsID:    NotificationNewTransactions{},
	NotificationRequestGetObjectsID:  NotificationRequestGetObjects{},
	NotificationResponseGetObjectsID: NotificationResponseGetObjects{},
	NotificationRequestChainID:       NotificationRequestChain{},
	NotificationResponseChainEntryID: NotificationResponseChainEntry{},
	NotificationTxPoolID:             NotificationTxPool{},
	NotificationNewLiteBlockID:       NotificationNewLiteBlock{},
}

func parseNotification(lc *LevinCommand) (interface{}, error) {
	if n, ok := mapNotificationID[lc.Command]; ok {
		notification := reflect.New(reflect.TypeOf(n))

		if err := binary.Unmarshal(lc.Payload, notification.Interface()); err != nil {
			return nil, err
		}

		return notification.Elem().Interface(), nil
	}

	return nil, errors.New(fmt.Sprintf("unknown notification ID: %d", lc.Command))
}

func (rb *RawBlock) ToBlock() (*cryptonote.Block, error) {
	block := &cryptonote.Block{}
	rawBlockReader := bytes.NewReader(rb.Block)
	if err := block.Deserialize(rawBlockReader); err != nil {
		return nil, err
	}

	return block, nil
}
