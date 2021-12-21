package cryptonote

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/r3volut1oner/go-karbo/config"
	"unsafe"
)

type ParentBlock struct {
	MajorVersion byte
	MinorVersion byte
	Timestamp    uint64
	Nonce        uint32
	Prev         Hash

	TransactionsCount     uint16
	BaseTransactionBranch HashList
	BaseTransaction       BaseTransaction
	BlockchainBranch      []Hash
}

type BlockHeader struct {
	MajorVersion      byte
	MinorVersion      byte
	Nonce             uint32
	Timestamp         uint64
	PreviousBlockHash Hash
}

type Block struct {
	Parent              ParentBlock
	CoinbaseTransaction Transaction
	TransactionsHashes  []Hash

	hash             *Hash
	hashTransactions *Hash

	BlockHeader
}

func (b *Block) Hash() *Hash {
	if b.hash == nil {
		var allBytesBuffer bytes.Buffer

		/**
		 * Write block header bytes
		 */
		allBytesBuffer.Write(b.BlockHeader.serialize())

		/**
		 * Write merkle root hash bytes
		 */
		baseTransactionHash := b.CoinbaseTransaction.Hash()
		hl := HashList{*baseTransactionHash}
		hl = append(hl, b.TransactionsHashes...)
		allBytesBuffer.Write(hl.merkleRootHash()[:])

		/**
		 * Write transactions number
		 */
		transactionCount := make([]byte, binary.MaxVarintLen64)
		written := binary.PutUvarint(transactionCount, uint64(len(hl)))
		allBytesBuffer.Write(transactionCount[:written])

		if b.MajorVersion == config.BlockMajorVersion2 || b.MajorVersion == config.BlockMajorVersion3 {
			allBytesBuffer.Write(b.Parent.serialize(true))
		}

		/**
		 * Create final hash bytes, by appending hash bytes length and the hash bytes
		 */
		allBytes := allBytesBuffer.Bytes()
		allBytesCount := make([]byte, binary.MaxVarintLen64)
		written = binary.PutUvarint(allBytesCount, uint64(len(allBytes)))

		var h bytes.Buffer
		h.Write(allBytesCount[:written])
		h.Write(allBytes)

		hashBytes := h.Bytes()
		b.hash = new(Hash)
		b.hash.FromBytes(hashBytes)
	}

	return b.hash
}

func (b *Block) Height() uint32 {
	if len(b.CoinbaseTransaction.Inputs) == 1 {
		i := b.CoinbaseTransaction.Inputs[0]

		if coinbase, ok := i.(InputCoinbase); ok {
			return coinbase.BlockIndex
		}
	}

	return 0
}

func (b *Block) Deserialize(r *bytes.Reader) error {
	if err := b.BlockHeader.deserialize(r); err != nil {
		return err
	}

	majorVersion := b.BlockHeader.MajorVersion
	if majorVersion == config.BlockMajorVersion2 || majorVersion == config.BlockMajorVersion3 {
		if err := b.Parent.deserialize(r); err != nil {
			return err
		}
	}

	if err := b.CoinbaseTransaction.Deserialize(r); err != nil {
		return err
	}

	hashesCount, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}

	hl := HashList{}
	for i := uint64(0); i < hashesCount; i++ {
		var h Hash
		if err := h.Read(r); err != nil {
			return err
		}
		hl = append(hl, h)
	}
	b.TransactionsHashes = hl

	return nil
}

func (b *Block) Serialize() []byte {
	var serialized bytes.Buffer

	serialized.Write(b.BlockHeader.serialize())

	mv := b.BlockHeader.MajorVersion
	if mv == config.BlockMajorVersion2 || mv == config.BlockMajorVersion3 {
		serialized.Write(b.Parent.serialize(false))
	}

	serialized.Write(b.CoinbaseTransaction.Serialize())

	buf := make([]byte, binary.MaxVarintLen64)
	written := binary.PutUvarint(buf, uint64(len(b.TransactionsHashes)))
	serialized.Write(buf[:written])

	for i := 0; i < len(b.TransactionsHashes); i++ {
		serialized.Write(b.TransactionsHashes[i][:])
	}

	return serialized.Bytes()
}

func (h *BlockHeader) deserialize(reader *bytes.Reader) error {
	var prev Hash
	var ts uint64
	var nonce uint32

	majorVersion, err := binary.ReadUvarint(reader)
	if err != nil {
		return nil
	}

	minorVersion, err := binary.ReadUvarint(reader)
	if err != nil {
		return nil
	}

	switch uint8(majorVersion) {
	case config.BlockMajorVersion2, config.BlockMajorVersion3:
		if err := prev.Read(reader); err != nil {
			return err
		}
	case config.BlockMajorVersion1, config.BlockMajorVersion4:
		ts, err = binary.ReadUvarint(reader)
		if err != nil {
			return err
		}

		if err := prev.Read(reader); err != nil {
			return err
		}

		if err := binary.Read(reader, binary.LittleEndian, &nonce); err != nil {
			return err
		}
	default:
		return errors.New("wrong block major version")
	}

	h.MajorVersion = uint8(majorVersion)
	h.MinorVersion = uint8(minorVersion)
	h.PreviousBlockHash = prev

	if ts != 0 {
		h.Timestamp = ts
	}

	if nonce != 0 {
		h.Nonce = nonce
	}

	return nil
}

func (h *BlockHeader) serialize() []byte {
	var serialized bytes.Buffer

	buf := make([]byte, binary.MaxVarintLen64)
	written := binary.PutUvarint(buf, uint64(h.MajorVersion))
	serialized.Write(buf[:written])

	written = binary.PutUvarint(buf, uint64(h.MinorVersion))
	serialized.Write(buf[:written])

	switch h.MajorVersion {
	case config.BlockMajorVersion2, config.BlockMajorVersion3:
		serialized.Write(h.PreviousBlockHash[:])
	case config.BlockMajorVersion1, config.BlockMajorVersion4:
		written = binary.PutUvarint(buf, h.Timestamp)
		serialized.Write(buf[:written])
		serialized.Write(h.PreviousBlockHash[:])
		_ = binary.Write(&serialized, binary.LittleEndian, h.Nonce)
	default:
	}

	return serialized.Bytes()
}

func (pb *ParentBlock) serialize(hashing bool) []byte {
	buf := make([]byte, binary.MaxVarintLen64)

	var serialized bytes.Buffer

	written := binary.PutUvarint(buf, uint64(pb.MajorVersion))
	serialized.Write(buf[:written])

	written = binary.PutUvarint(buf, uint64(pb.MinorVersion))
	serialized.Write(buf[:written])

	written = binary.PutUvarint(buf, pb.Timestamp)
	serialized.Write(buf[:written])

	serialized.Write(pb.Prev[:])

	_ = binary.Write(&serialized, binary.LittleEndian, pb.Nonce)

	if hashing {
		th := pb.BaseTransaction.Hash()
		h := pb.BaseTransactionBranch.TreeHashFromBranch(*th)

		serialized.Write(h[:])
	}

	written = binary.PutUvarint(buf, uint64(pb.TransactionsCount))
	serialized.Write(buf[:written])

	for _, tb := range pb.BaseTransactionBranch {
		serialized.Write(tb[:])
	}

	serialized.Write(pb.BaseTransaction.serialize())

	for _, h := range pb.BlockchainBranch {
		serialized.Write(h[:])
	}

	return serialized.Bytes()
}

func (pb *ParentBlock) deserialize(r *bytes.Reader) error {
	var prev Hash
	var nonce uint32

	majorVersion, err := binary.ReadUvarint(r)
	if err != nil {
		return nil
	}

	minorVersion, err := binary.ReadUvarint(r)
	if err != nil {
		return nil
	}

	timestamp, err := binary.ReadUvarint(r)
	if err != nil {
		return nil
	}

	if err := prev.Read(r); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &nonce); err != nil {
		return err
	}

	txCount, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}

	var baseTxBranch HashList
	branchSize := treeDepth(int(txCount))
	for i := 0; i < branchSize; i++ {
		var th Hash
		if err := th.Read(r); err != nil {
			return err
		}
		baseTxBranch = append(baseTxBranch, th)
	}

	baseTx := BaseTransaction{&TransactionPrefix{}, nil}
	if err := baseTx.deserialize(r); err != nil {
		return err
	}

	tef, err := baseTx.ParseExtra()
	if err != nil {
		return err
	}

	if tef.MiningTag == nil {
		return errors.New("can't get extra merge mining tag")
	}

	if tef.MiningTag.Depth > 8*32 {
		return errors.New("wrong merge mining tag depth")
	}

	var blockchainBranch HashList
	for i := uint64(0); i < tef.MiningTag.Depth; i++ {
		var h Hash
		if err := h.Read(r); err != nil {
			return err
		}
		blockchainBranch = append(blockchainBranch, h)
	}

	pb.MajorVersion = byte(majorVersion)
	pb.MinorVersion = byte(minorVersion)
	pb.Timestamp = timestamp
	pb.Nonce = nonce
	pb.Prev = prev

	pb.TransactionsCount = uint16(txCount)
	pb.BaseTransactionBranch = baseTxBranch
	pb.BaseTransaction = baseTx
	pb.BlockchainBranch = blockchainBranch

	return nil
}

//size_t tree_depth(size_t count) {
//	size_t i;
//	size_t depth = 0;
//	assert(count > 0);
//	for (i = sizeof(size_t) << 2; i > 0; i >>= 1) {
//		if (count >> i > 0) {
//			count >>= i;
//			depth += i;
//		}
//	}
//	return depth;
//}
func treeDepth(count int) int {
	depth := 0

	for i := int(unsafe.Sizeof(count)) << 2; i > 0; i >>= 1 {
		if count>>i > 0 {
			count >>= i
			depth += i
		}
	}

	return depth
}
