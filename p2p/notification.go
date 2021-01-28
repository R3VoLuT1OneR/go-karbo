package p2p

import "github.com/r3volut1oner/go-karbo/cryptonote"

type NotificationTxPool struct {
	Transactions []cryptonote.Hash `binary:"txs"`
}
