package config

import "C"
import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/r3volut1oner/go-karbo/utils"
	"math"
	"sort"
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

	allowLowDifficulty bool
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
	testnet.allowLowDifficulty = true

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

func (n *Network) MaxTransactionSize(height uint32) uint64 {
	if height > UpgradeHeightV4 {
		return transactionMaxSize
	}

	// Return maximum unachievable possible transaction size by default
	return math.MaxUint64
}

func (n *Network) GetBlockMajorVersionForHeight(h uint32) byte {
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

func (n *Network) DifficultyBlocksCountByBlockVersion(majorVersion byte) int {
	if majorVersion >= BlockMajorVersion5 {
		return difficultyWindow4 + 1
	}

	if majorVersion == BlockMajorVersion3 || majorVersion == BlockMajorVersion4 {
		return difficultyWindow3 + 1
	}

	if majorVersion == BlockMajorVersion2 {
		return difficultyWindow2
	}

	return difficultyWindow + difficultyLag
}

func (n *Network) NextDifficulty(h uint32, nextBlockMajorVersion byte, timestamps []uint64, cumulativeDifficulties []uint64) (uint64, error) {
	if nextBlockMajorVersion >= BlockMajorVersion5 {
		return n.nextDifficultyV5(h, nextBlockMajorVersion, timestamps, cumulativeDifficulties)
	}

	if nextBlockMajorVersion == BlockMajorVersion4 {
		return n.nextDifficultyV4(h, timestamps, cumulativeDifficulties)
	}

	if nextBlockMajorVersion == BlockMajorVersion3 {
		return n.nextDifficultyV3(timestamps, cumulativeDifficulties)
	}

	if nextBlockMajorVersion == BlockMajorVersion2 {
		return n.nextDifficultyV2(timestamps, cumulativeDifficulties)
	}

	return n.nextDifficultyV1(timestamps, cumulativeDifficulties)
}

func (n *Network) MinimalFee(height uint32) uint64 {
	if height <= UpgradeHeightV3s1 {
		return minimumFeeV1
	}

	if height > UpgradeHeightV3s1 && height <= UpgradeHeightV4 {
		return minimumFeeV2
	}

	if height > UpgradeHeightV4 && height < UpgradeHeightV4s3 {
		return minimumFeeV1
	}

	return minimumFeeV3
}

func (n *Network) nextDifficultyV5(height uint32, majorVersion byte, timestamps []uint64, cumulativeDifficulties []uint64) (uint64, error) {
	// LWMA-1 difficulty algorithm
	// Copyright (c) 2017-2018 Zawy, MIT License
	// See commented link below for required config file changes. Fix FTL and MTP.
	// https://github.com/zawy12/difficulty-algorithms/issues/3

	// begin reset difficulty for new epoch
	if height == UpgradeHeightV5 {
		return cumulativeDifficulties[0] / uint64(height) / resetWorkFactorV5, nil
	}

	count := uint32(n.DifficultyBlocksCountByBlockVersion(majorVersion) - 1)
	if height > UpgradeHeightV5 && height < UpgradeHeightV5+count {
		offset := count - (height - UpgradeHeightV5)
		timestamps = timestamps[offset:]
		cumulativeDifficulties = cumulativeDifficulties[offset:]
	}

	// end reset difficulty for new epoch
	if len(timestamps) != len(cumulativeDifficulties) {
		return 0, errors.New(fmt.Sprintf(
			"assertion failed - timestamp count (%d) must match cumulativeDifficulties count (%d)",
			len(timestamps),
			len(cumulativeDifficulties),
		))
	}

	T := int64(difficultyTarget)
	// adjust for new epoch difficulty reset, N should be by 1 block smaller
	N := utils.MinInt64(difficultyWindow4, int64(len(cumulativeDifficulties)-1))
	L := int64(0)

	thisTimestamp := int64(0)
	previousTimestamp := int64(timestamps[0]) - T
	for i := int64(1); i <= N; i++ {
		ft := int64(timestamps[i])

		// Safely prevent out-of-sequence timestamps
		if ft > previousTimestamp {
			thisTimestamp = ft
		} else {
			thisTimestamp = previousTimestamp + 1
		}

		L += i * utils.MinInt64(6*T, thisTimestamp-previousTimestamp)
		previousTimestamp = thisTimestamp
	}

	if L < N*N*T/20 {
		L = N * N * T / 20
	}

	nextD := int64(0)
	avgD := int64(cumulativeDifficulties[N]) - int64(cumulativeDifficulties[0])/N

	// Prevent round off error for small D and overflow for large D.
	if avgD > 2000000*N*N*T {
		nextD = (avgD / (200 * L)) * (N * (N + 1) * T * 99)
	} else {
		nextD = (avgD * N * (N + 1) * T * 99) / (L * 200)
	}

	// Optional. Make all insignificant digits zero for easy reading.
	i := int64(1000000000)
	for i > 1 {
		if nextD > i*100 {
			nextD = (((nextD + i) / 2) / i) * i
			break
		}

		i /= 10
	}

	if !n.allowLowDifficulty && nextD < minimalDifficulty {
		nextD = minimalDifficulty
	}

	return uint64(nextD), nil
}

func (n *Network) nextDifficultyV4(height uint32, timestamps []uint64, cumulativeDifficulties []uint64) (uint64, error) {
	// LWMA-2 / LWMA-3 difficulty algorithm
	// Copyright (c) 2017-2018 Zawy, MIT License
	// https://github.com/zawy12/difficulty-algorithms/issues/3
	// with modifications by Ryo Currency developers

	T := int64(difficultyTarget)
	N := int64(difficultyWindow3)

	var L, ST, sum3ST int64 = 0, 0, 0

	if len(timestamps) != len(cumulativeDifficulties) {
		return 0, errors.New(fmt.Sprintf(
			"assertion failed - timestamp count (%d) must match cumulativeDifficulties count (%d)",
			len(timestamps),
			len(cumulativeDifficulties),
		))
	}

	if len(timestamps) > int(N+1) {
		return 0, errors.New(fmt.Sprintf(
			"assertion failed - timestamps size (%d) bigger than (%d)",
			len(timestamps),
			N+1,
		))
	}

	var maxTS, prevMaxTS int64 = 0, int64(timestamps[0])
	var lwma3Height = UpgradeHeightV4s1

	for i := int64(1); i <= N; i++ {
		tf := int64(timestamps[i])
		tl := int64(timestamps[i-1])

		if height < lwma3Height {
			ST = utils.ClampInt64(-6*T, tf-tl, 6*T)
		} else { // LWMA-3
			if tf > prevMaxTS {
				maxTS = tf
			} else {
				maxTS = prevMaxTS + 1
			}

			ST = utils.MinInt64(6*T, maxTS-prevMaxTS)
			prevMaxTS = maxTS
		}

		L += ST * i
		if i > N-3 {
			sum3ST += ST
		}
	}

	nextD := int64(cumulativeDifficulties[N]-cumulativeDifficulties[0]) * T * (N * 1) / 2 * L
	nextD = (nextD * 99) / 100

	prevD := int64(cumulativeDifficulties[N] - cumulativeDifficulties[N-1])
	nextD = utils.ClampInt64(prevD*67/100, nextD, prevD*150/100)

	if sum3ST < (8*T)/10 {
		nextD = (prevD * 110) / 100
	}

	if !n.allowLowDifficulty && nextD < minimalDifficulty {
		nextD = minimalDifficulty
	}

	return uint64(nextD), nil
}

func (n *Network) nextDifficultyV3(timestamps []uint64, cumulativeDifficulties []uint64) (uint64, error) {
	// LWMA difficulty algorithm
	// Copyright (c) 2017-2018 Zawy
	// MIT license http://www.opensource.org/licenses/mit-license.php.
	// This is an improved version of Tom Harding's (Deger8) "WT-144"
	// Karbowanec, Masari, Bitcoin Gold, and Bitcoin Cash have contributed.
	// See https://github.com/zawy12/difficulty-algorithms/issues/1 for other algos.
	// Do not use "if solvetime < 0 then solvetime = 1" which allows a catastrophic exploit.
	// T= target_solvetime;
	// N = int(45 * (600 / T) ^ 0.3));

	T := int64(difficultyTarget)
	N := difficultyWindow3

	// return a difficulty of 1 for first 3 blocks if it's the start of the chain
	if len(timestamps) < 4 {
		return 1, nil
	}

	// otherwise, use a smaller N if the start of the chain is less than N+1
	if len(timestamps) < N+1 {
		N = len(timestamps) - 1
	}

	if len(timestamps) > N+1 {
		timestamps = timestamps[N-1:]
		cumulativeDifficulties = cumulativeDifficulties[N-1:]
	}

	// To get an average solvetime to within +/- ~0.1%, use an adjustment factor.
	adjust := 0.998

	// The divisor k normalizes LWMA.
	k := float64(N * (N + 1) / 2.0)

	LWMA := float64(0)
	sumInverseD := float64(0)

	solveTime := int64(0)
	difficulty := int64(0)

	// Loop through N most recent blocks.
	for i := 1; i <= N; i++ {
		solveTime = int64(timestamps[i]) - int64(timestamps[i-1])
		solveTime = utils.MinInt64(T*7, utils.MaxInt64(solveTime, -6*T))
		difficulty = int64(cumulativeDifficulties[i]) - int64(cumulativeDifficulties[i-1])
		LWMA += float64(solveTime*1) / k
		sumInverseD += 1 / float64(difficulty)
	}

	// Keep LWMA sane in case something unforeseen occurs.
	if int64(math.Round(LWMA)) < T/20 {
		LWMA = float64(T) / 20
	}

	harmonicMeanD := float64(N) / sumInverseD * adjust
	nextDifficulty := uint64(harmonicMeanD * float64(T) / LWMA)

	if !n.allowLowDifficulty && nextDifficulty < minimalDifficulty {
		nextDifficulty = minimalDifficulty
	}

	return nextDifficulty, nil
}

func (n *Network) nextDifficultyV2(timestamps []uint64, cumulativeDifficulties []uint64) (uint64, error) {
	// Difficulty calculation v. 2
	// based on Zawy difficulty algorithm v1.0
	// next Diff = Avg past N Diff * TargetInterval / Avg past N solve times
	// as described at https://github.com/monero-project/research-lab/issues/3
	// Window time span and total difficulty is taken instead of average as suggested by Nuclear_chaos

	if len(timestamps) != len(cumulativeDifficulties) {
		return 0, errors.New(fmt.Sprintf(
			"assertion failed - timestamp count (%d) must match cumulativeDifficulties count (%d)",
			len(timestamps),
			len(cumulativeDifficulties),
		))
	}

	if len(timestamps) > difficultyWindow2 {
		timestamps = timestamps[:difficultyWindow2]
		cumulativeDifficulties = cumulativeDifficulties[:difficultyWindow2]
	}

	length := len(timestamps)
	if length <= 1 {
		return 1, nil
	}

	sort.Slice(timestamps, func(i, j int) bool { return timestamps[i] < timestamps[j] })

	timeSpan := timestamps[len(timestamps)-1] - timestamps[0]
	if timeSpan == 0 {
		timeSpan = 1
	}

	totalWork := cumulativeDifficulties[len(cumulativeDifficulties)-1] - cumulativeDifficulties[0]
	if totalWork <= 0 {
		return 0, errors.New(fmt.Sprintf("assertion failed - total work %d <= 0", totalWork))
	}

	low, high := utils.Mul128(totalWork, difficultyTarget)
	if high != 0 {
		return 0, nil
	}

	nextDiffZ := low / timeSpan

	if !n.allowLowDifficulty && nextDiffZ < minimalDifficulty {
		nextDiffZ = minimalDifficulty
	}

	return nextDiffZ, nil
}

func (n *Network) nextDifficultyV1(timestamps []uint64, cumulativeDifficulties []uint64) (uint64, error) {
	if len(timestamps) != len(cumulativeDifficulties) {
		return 0, errors.New(fmt.Sprintf(
			"assertion failed - timestamp count (%d) must match cumulativeDifficulties count (%d)",
			len(timestamps),
			len(cumulativeDifficulties),
		))
	}

	if len(timestamps) > difficultyWindow {
		timestamps = timestamps[:difficultyWindow]
		cumulativeDifficulties = cumulativeDifficulties[:difficultyWindow]
	}

	length := len(timestamps)
	if length <= 1 {
		return 1, nil
	}

	sort.Slice(timestamps, func(i, j int) bool { return timestamps[i] < timestamps[j] })

	var cutBegin, cutEnd int

	//if !(2*difficultyCut <= difficultyWindow-2) {
	//	return 0, errors.New("assertion failed")
	//}

	if length <= difficultyWindow-(2*difficultyCut) {
		cutBegin = 0
		cutEnd = length
	} else {
		cutBegin = (length - (difficultyWindow - (2 * difficultyCut)) + 1) / 2
		cutEnd = cutBegin + (difficultyWindow - (2 * difficultyCut))
	}

	if !(cutBegin+2 <= cutEnd && (cutEnd) <= length) {
		return 0, errors.New("assertion failed")
	}

	timeSpan := timestamps[cutEnd-1] - timestamps[cutBegin]
	if timeSpan == 0 {
		timeSpan = 1
	}

	totalWork := cumulativeDifficulties[cutEnd-1] - cumulativeDifficulties[cutBegin]
	if totalWork <= 0 {
		return 0, errors.New(fmt.Sprintf("assertion failed - total work %d <= 0", totalWork))
	}

	low, high := utils.Mul128(totalWork, difficultyTarget)

	if high != 0 || math.MaxUint64-low < timeSpan-1 {
		return 0, nil
	}

	return (low + timeSpan - 1) / timeSpan, nil
}
