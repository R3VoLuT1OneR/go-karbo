package config

const (
	DifficultyTarget = uint64(240)

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

	MaxBlockNumber   = uint64(500000000)
	MaxBlockBlobSize = uint64(500000000)
	MaxTxSize        = uint64(1000000000)

	UpgradeHeightV2   = uint32(60000)
	UpgradeHeightV3   = uint32(216000)
	UpgradeHeightV3s1 = uint32(216394)
	UpgradeHeightV4   = uint32(266000)
	UpgradeHeightV4s1 = uint32(300000)
	UpgradeHeightV4s2 = uint32(500000)
	UpgradeHeightV5   = uint32(4294967294)

	blockFutureTimeLimit   = DifficultyTarget * 7
	blockFutureTimeLimitV1 = DifficultyTarget * 3

	blockTimestampCheckWindow   = 60
	blockTimestampCheckWindowV1 = 11

	minedMoneyUnlockWindow = 10

	MaxExtraSize = 1024
)
