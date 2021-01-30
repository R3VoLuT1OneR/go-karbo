package p2p

import (
	"errors"
	"fmt"
)

type PeerStore struct {
	peers map[uint64]*Peer
}

func NewPeerStore() *PeerStore {
	return &PeerStore{
		peers: map[uint64]*Peer{},
	} 
}

func (ps *PeerStore) Add(p *Peer) error {
	if p.ID == 0 {
		return errors.New("peer must have ID for save it in store")
	}

	if _, ok := ps.peers[p.ID]; ok {
		return errors.New(fmt.Sprintf("peer '%v' already in the store", p.ID))
	}

	ps.peers[p.ID] = p

	return nil
}

func (ps *PeerStore) Remove(p *Peer) error {
	if _, ok := ps.peers[p.ID]; !ok {
		return errors.New("peer id not found in store")
	}

	delete(ps.peers, p.ID)
	return nil
}

func (ps *PeerStore) Has(ID uint64) bool {
	_, ok := ps.peers[ID]

	return ok
}

func (ps *PeerStore) Get(ID uint64) *Peer {
	return ps.peers[ID]
}
