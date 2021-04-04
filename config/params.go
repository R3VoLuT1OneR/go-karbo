package config

import "github.com/google/uuid"

const (
	TransactionVersion1 byte = 1

	P2PVersion1 uint8 = 1
	P2PVersion2 uint8 = 2
	P2PVersion3 uint8 = 3
	P2PVersion4 uint8 = 4

	BlockMajorVersion1 byte = 1
	BlockMajorVersion2 byte = 2
	BlockMajorVersion3 byte = 3
	BlockMajorVersion4 byte = 4
	BlockMajorVersion5 byte = 5

	BlockMinorVersion0 byte = 0
	BlockMinorVersion1 byte = 1

	MaxBlockNumber 			= uint64(500000000)
	MaxBlockBlobSize 		= uint64(500000000)
	MaxTxSize 				= uint64(1000000000)
)

// Network represents network params
type Network struct {

	NetworkID uuid.UUID

	MaxBlockNumber uint64
	MaxBlockBlobSize uint64
	MaxTxSize uint64

	Name string
	Ticker string
	GenesisCoinbaseTxHex string

	// PublicAddressBase58Prefix address prefix
	PublicAddressBase58Prefix uint64

	TxProofBase58Prefix uint64

	ReserveProofBase58Prefix uint64

	KeysSignatureBase58Prefix uint64

	P2PMinimumVersion byte
	P2PCurrentVersion byte

	CurrentTransactionVersion byte

	// SeadNodes List of basic sead nodes
	SeedNodes []string
}
