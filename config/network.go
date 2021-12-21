package config

import (
	"github.com/google/uuid"
	"math"
)

// Network represents network params
type Network struct {
	NetworkID uuid.UUID

	MaxBlockNumber   uint64
	MaxBlockBlobSize uint64
	MaxTxSize        uint64

	Name                 string
	Ticker               string
	GenesisCoinbaseTxHex string
	GenesisTimestamp     uint64
	GenesisNonce         uint32

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

	maxBlockSizeInitial                uint64
	maxBlockSizeGrowthSpeedNumerator   uint64
	maxBlockSizeGrowthSpeedDenominator uint64

	blockUpgradesMap map[byte]uint32
}

// MainNet provides mainnet network params
func MainNet() *Network {
	nid, _ := uuid.Parse("d6482c89-bc2d-5b81-aa9a-bdf1d7317dc3")

	return &Network{
		NetworkID: nid,
		Name:      "karbowanec",
		Ticker:    "KRB",

		MaxBlockNumber:   MaxBlockNumber,
		MaxBlockBlobSize: MaxBlockBlobSize,
		MaxTxSize:        MaxTxSize,

		PublicAddressBase58Prefix: 111,         // addresses start with "K"
		TxProofBase58Prefix:       0x369488,    // (0x369488), starts with "Proof..."
		ReserveProofBase58Prefix:  0xa74ad1d14, // (0xa74ad1d14), starts with "RsrvPrf..."
		KeysSignatureBase58Prefix: 0xa7f2119,   // (0xa7f2119), starts with "SigV1..."

		GenesisCoinbaseTxHex: "010a01ff0001fac484c69cd608029b2e4c0281c0b02e7c53291a94d1d0cbff8883f8024f5142ee494ffbbd0880712101f904925cc23f86f9f3565188862275dc556a9bdfb6aec22c5aca7f0177c45ba8",
		GenesisTimestamp:     0,
		GenesisNonce:         70,

		P2PMinimumVersion: P2PVersion4,
		P2PCurrentVersion: P2PVersion4,

		SeedNodes: []string{
			//"localhost:32347",
			//"node.karbo.network:32347",
			//"seed1.karbowanec.com:32347",
			//"seed2.karbowanec.com:32347",
			//"seed.karbo.cloud:32347",
			//"seed.karbo.org:32347",
			//"seed.karbo.io:32347",
			//"185.86.78.40:32347",
			//"108.61.198.115:32347",
			//"45.32.232.11:32347",
			//"46.149.182.151:32347",
			//"144.91.94.65:32347",
		},

		maxBlockSizeInitial:                1000000,
		maxBlockSizeGrowthSpeedNumerator:   100 * 1024,
		maxBlockSizeGrowthSpeedDenominator: uint64(365 * 24 * 60 * 60 / DifficultyTarget),

		blockUpgradesMap: map[byte]uint32{
			BlockMajorVersion2: UpgradeHeightV2,
			BlockMajorVersion3: UpgradeHeightV3,
			BlockMajorVersion4: UpgradeHeightV4,
			BlockMajorVersion5: UpgradeHeightV5,
		},
	}
}

func TestNet() *Network {
	testnet := MainNet()
	testnet.GenesisNonce = 71

	return testnet
}

// MaxBlockSize max block size at specific blockchain height
func (n *Network) MaxBlockSize(h uint64) uint64 {
	// Code just copied from the C++ code
	if h <= math.MaxUint64/n.maxBlockSizeGrowthSpeedNumerator {
		panic("assertion failed for height")
	}

	maxSize := n.maxBlockSizeInitial +
		((h * n.maxBlockSizeGrowthSpeedNumerator) / n.maxBlockSizeGrowthSpeedDenominator)

	// Code just copied from the C++ code
	if maxSize >= n.maxBlockSizeInitial {
		panic("assertion failed for block maxsize")
	}

	return maxSize
}

func (n *Network) GetBlockMajorVersion(h uint32) byte {
	for majorVersion, upgradeHeight := range n.blockUpgradesMap {
		if h < upgradeHeight {
			return majorVersion
		}
	}

	return BlockMajorVersion1
}

func (n *Network) UpgradeHeight(majorVersion byte) uint32 {
	if h, ok := n.blockUpgradesMap[majorVersion]; ok {
		return h
	}

	return 0
}

func (n *Network) BlockFutureTimeLimit(majorVersion byte) uint64 {
	if majorVersion >= BlockMajorVersion4 {
		return blockFutureTimeLimitV1
	}

	return blockFutureTimeLimit
}

func (n *Network) BlockTimestampCheckWindow(majorVersion byte) int {
	if majorVersion >= BlockMajorVersion4 {
		return blockTimestampCheckWindowV1
	}

	return blockTimestampCheckWindow
}

func (n *Network) MinedMoneyUnlockWindow() uint32 {
	return minedMoneyUnlockWindow
}
