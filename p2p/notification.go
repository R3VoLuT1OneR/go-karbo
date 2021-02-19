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

type RawBlock struct {
	Block 			[]byte 		`binary:"block"`
	Transactions 	[][]byte 	`binary:"txs,array,omitempty"`
}

type NotificationNewBlock struct {
	Block 					RawBlock	`binary:"b"`
	CurrentBlockchainHeight uint32 		`binary:"current_blockchain_height"`
	Hop             		uint32 	 	`binary:"hop"`
}

type NotificationNewTransactions struct {
	Stem         bool              `binary:"stem"`
	Transactions []cryptonote.Hash `binary:"txs,binary"`
}

type NotificationRequestGetObjects struct {
	Transactions []cryptonote.Hash `binary:"txs,binary"`
	Blocks       []cryptonote.Hash `binary:"blocks,binary"`
}

type NotificationResponseGetObjects struct {
	// Exists in old legacy code in definition but not exists in notification
	Transactions 			[]string 			`binary:"txs,omitempty"`
	Blocks                  []RawBlock          `binary:"blocks,array"`
	CurrentBlockchainHeight uint32 				`binary:"current_blockchain_height"`
	MissedIds 				[]cryptonote.Hash 	`binary:"missed_ids,binary,omitempty"`
}

type NotificationTxPool struct {
	Transactions []cryptonote.Hash `binary:"txs"`
}

type NotificationNewLiteBlock struct {
	CurrentBlockchainHeight uint32 	 	`binary:"current_blockchain_height"`
	Hop                     uint32 	 	`binary:"hop"`
	Block 					[]byte 		`binary:"block"`
}

type NotificationRequestChain struct {
	Blocks []cryptonote.Hash `binary:"block_ids,binary"`
}

type NotificationResponseChainEntry struct {
	Start    uint32            `binary:"start_height"`
	Total    uint32            `binary:"total_height"`
	BlockIds []cryptonote.Hash `binary:"m_block_ids,binary"`
}

var mapNotificationID = map[uint32]interface{}{
	NotificationNewBlockID: 			NotificationNewBlock{},
	NotificationNewTransactionsID: 		NotificationNewTransactions{},
	NotificationRequestGetObjectsID:    NotificationRequestGetObjects{},
	NotificationResponseGetObjectsID:   NotificationResponseGetObjects{},
	NotificationTxPoolID: 				NotificationTxPool{},
	NotificationNewLiteBlockID: 		NotificationNewLiteBlock{},
	NotificationRequestChainID: 		NotificationRequestChain{},
	NotificationResponseChainEntryID: 	NotificationResponseChainEntry{},
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