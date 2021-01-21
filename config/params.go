package config

import "github.com/google/uuid"

const (
	P2PVersion1 uint8 = 1
	P2PVersion2 uint8 = 2
	P2PVersion3 uint8 = 3
	P2PVersion4 uint8 = 4
)

// Params represents network params
type Params struct {

	NetworkID uuid.UUID

	// PublicAddressBase58Prefix address prefix
	PublicAddressBase58Prefix uint64

	TxProofBase58Prefix uint64

	ReserveProofBase58Prefix uint64

	KeysSignatureBase58Prefix uint64

	P2PMinimumVersion uint8
	P2PCurrentVersion uint8

	// SeadNodes List of basic sead nodes
	SeedNodes []string
}
