package p2p

import (
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/cryptonote"
	"github.com/r3volut1oner/go-karbo/encoding/binary"
	"reflect"
)

const (
	NotificationBoolBase 				= 2000
	NotificationNewBlockID 				= NotificationBoolBase + 1
	NotificationNewTransactionsID 		= NotificationBoolBase + 2
	NotificationRequestGetObjectsID 	= NotificationBoolBase + 3
	NotificationResponseGetObjectsID 	= NotificationBoolBase + 4
	NotificationRequestChainID 			= NotificationBoolBase + 6
	NotificationResponseChainEntryID	= NotificationBoolBase + 7
	NotificationTxPoolID 				= NotificationBoolBase + 8
	NotificationNewLiteBlockID 			= NotificationBoolBase + 9
	NotificationMissingTxsID 			= NotificationBoolBase + 10
)

type NotificationNewBlock struct {

}

type NotificationNewTransactions struct {
	Transactions []cryptonote.Hash `binary:"txs"`
}

type NotificationTxPool struct {
	Transactions []cryptonote.Hash `binary:"txs"`
}

type NotificationNewLiteBlock struct {
	Block 					[]byte `binary:"block"`
	CurrentBlockchainHeight uint32 `binary:"current_blockchain_height"`
	Hop                     uint32 `binary:"hop"`
}

type NotificationRequestChain struct {
	Blocks []cryptonote.Hash `binary:"block_ids"`
}

var mapNotificationID = map[uint32]interface{}{
	NotificationNewBlockID: 		NotificationNewBlock{},
	NotificationNewTransactionsID: 	NotificationNewTransactions{},
	NotificationTxPoolID: 			NotificationTxPool{},
	NotificationNewLiteBlockID: 	NotificationNewLiteBlock{},
	NotificationRequestChainID: 	NotificationRequestChain{},
}

func parseNotification(lc *LevinCommand) (interface{}, error) {
	if n, ok := mapNotificationID[lc.Command]; ok {
		notification := reflect.New(reflect.TypeOf(n))

		if err := binary.Unmarshal(lc.Payload, notification.Interface()); err != nil {
			return nil, err
		}

		return notification.Elem().Interface(), nil
	}

	return nil, errors.New(fmt.Sprintf("unknown command ID: %d", lc.Command))
}

func newRequestChain(n *config.Network) (*NotificationRequestChain, error) {
	topBlock, err := cryptonote.GenerateGenesisBlock(n)
	if err != nil {
		return nil, err
	}

	hash, err := topBlock.Hash()
	if err != nil {
		return nil, err
	}

	return &NotificationRequestChain{[]cryptonote.Hash{*hash}}, nil
}