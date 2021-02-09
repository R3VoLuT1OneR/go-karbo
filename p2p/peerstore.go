package p2p

import (
	"errors"
	"fmt"
)


type peerStore struct {
	white 	*peerList
	grey 	*peerList
}

type peerList struct {
	peers map[uint64]*Peer
}

func NewPeerStore() *peerStore {
	return &peerStore{
		white: &peerList{map[uint64]*Peer{}},
		grey: &peerList{map[uint64]*Peer{}},
	}
}

func (ps *peerStore) toWhite(p *Peer) error {
	_ = ps.grey.Remove(p)

	if err := ps.white.Add(p); err != nil {
		return err
	}

	return nil
}

func (ps *peerStore) toGrey(p *Peer) error {
	_ = ps.white.Remove(p)

	if err := ps.grey.Add(p); err != nil {
		return err
	}

	return nil
}

func (pl *peerList) Add(p *Peer) error {
	if p.ID == 0 {
		return errors.New("peer must have ID for save it in store")
	}

	if _, ok := pl.peers[p.ID]; ok {
		return errors.New(fmt.Sprintf("peer '%v' already in the store", p.ID))
	}

	pl.peers[p.ID] = p

	return nil
}

func (pl *peerList) Remove(p *Peer) error {
	if _, ok := pl.peers[p.ID]; !ok {
		return errors.New("peer id not found in store")
	}

	delete(pl.peers, p.ID)
	return nil
}

func (pl *peerList) Has(ID uint64) bool {
	_, ok := pl.peers[ID]

	return ok
}

func (pl *peerList) Get(ID uint64) *Peer {
	return pl.peers[ID]
}
