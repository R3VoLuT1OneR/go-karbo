package cryptonote

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
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
	BaseTransactionBranch []Hash
	BaseTransaction       BaseTransaction
	BlockchainBranch      []Hash
}

type BlockHeader struct {
	MajorVersion 		byte
	MinorVersion 		byte
	Nonce        		uint32
	Timestamp    		uint64
	Prev         		Hash
}

type Block struct {
	BlockHeader

	Parent 				ParentBlock
	Transaction 		Transaction
	TransactionsHashes 	[]Hash

	hash 				*Hash
	hashTransactions 	*Hash
}

func (b *Block) Hash() (*Hash, error) {
	if b.hash == nil {
		var allBytesBuffer bytes.Buffer

		/**
		 * Write block header bytes
		 */
		headerBytes, err := b.BlockHeader.serialize()
		if err != nil {
			return nil, err
		}
		allBytesBuffer.Write(headerBytes)

		/**
		 * Write merkle root hash bytes
		 */
		baseTransactionHash, err := b.Transaction.Hash()
		if err != nil {
			return nil, err
		}

		hl := HashList{*baseTransactionHash}
		hl = append(hl, b.TransactionsHashes...)

		mrHash, err := hl.merkleRootHash()
		if err != nil {
			return nil, err
		}

		allBytesBuffer.Write(mrHash[:])

		/**
		 * Write transactions number
		 */
		transactionCount := make([]byte, binary.MaxVarintLen64)
		written := binary.PutUvarint(transactionCount, uint64(len(hl)))
		allBytesBuffer.Write(transactionCount[:written])

		if b.MajorVersion == config.BlockMajorVersion2 || b.MajorVersion == config.BlockMajorVersion3 {
			bs, err := b.Parent.serialize(true)
			if err != nil {
				return nil, err
			}
			allBytesBuffer.Write(bs)
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
		b.hash.FromBytes(&hashBytes)
	}

	return b.hash, nil
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

	if err := b.Transaction.Deserialize(r); err != nil {
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

func (b *Block) Serialize() ([]byte, error) {
	var serialized bytes.Buffer

	hb, err := b.BlockHeader.serialize()
	if err != nil {
		return nil, err
	}
	serialized.Write(hb)

	mv := b.BlockHeader.MajorVersion
	if mv == config.BlockMajorVersion2 || mv == config.BlockMajorVersion3 {
		pbh, err := b.Parent.serialize(false)
		if err != nil {
			return nil, err
		}
		serialized.Write(pbh)
	}

	tb, err := b.Transaction.Serialize()
	if err != nil {
		return nil, err
	}
	serialized.Write(tb)

	buf := make([]byte, binary.MaxVarintLen64)
	written := binary.PutUvarint(buf, uint64(len(b.TransactionsHashes)))
	serialized.Write(buf[:written])

	for i := 0; i < len(b.TransactionsHashes); i++ {
		serialized.Write(b.TransactionsHashes[i][:])
	}

	return serialized.Bytes(), nil
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
	h.Prev = prev

	if ts != 0 {
		h.Timestamp = ts
	}

	if nonce != 0 {
		h.Nonce = nonce
	}

	return nil
}

func (h *BlockHeader) serialize() ([]byte, error) {
	var serialized bytes.Buffer

	buf := make([]byte, binary.MaxVarintLen64)
	written := binary.PutUvarint(buf, uint64(h.MajorVersion))
	serialized.Write(buf[:written])

	written = binary.PutUvarint(buf, uint64(h.MinorVersion))
	serialized.Write(buf[:written])

	switch h.MajorVersion {
	case config.BlockMajorVersion2, config.BlockMajorVersion3:
		serialized.Write(h.Prev[:])
	case config.BlockMajorVersion1, config.BlockMajorVersion4:
		written = binary.PutUvarint(buf, h.Timestamp)
		serialized.Write(buf[:written])
		serialized.Write(h.Prev[:])
		if err := binary.Write(&serialized, binary.LittleEndian, h.Nonce); err != nil {
			return nil, errors.New("failed to write block nonce")
		}
	default:
		return nil, errors.New("wrong block major version")
	}

	return serialized.Bytes(), nil
}

func (pb *ParentBlock) serialize(hashing bool) ([]byte, error) {
	buf := make([]byte, binary.MaxVarintLen64)

	var serialized bytes.Buffer

	written := binary.PutUvarint(buf, uint64(pb.MajorVersion))
	serialized.Write(buf[:written])

	written = binary.PutUvarint(buf, uint64(pb.MinorVersion))
	serialized.Write(buf[:written])

	written = binary.PutUvarint(buf, pb.Timestamp)
	serialized.Write(buf[:written])

	serialized.Write(pb.Prev[:])

	if err := binary.Write(&serialized, binary.LittleEndian, pb.Nonce); err != nil {
		return nil, err
	}

	if hashing {
		h, err := pb.BaseTransaction.Hash()
		if err != nil {
			return nil, err
		}
		hl := HashList{*h}

		tmh, err := hl.merkleRootHash()
		if err != nil {
			return nil, err
		}

		serialized.Write(tmh[:])
	}

	written = binary.PutUvarint(buf, uint64(pb.TransactionsCount))
	serialized.Write(buf[:written])

	for i := 0; i < len(pb.BaseTransactionBranch); i++ {
		serialized.Write(pb.BaseTransactionBranch[i][:])
	}

	btb, err := pb.BaseTransaction.serialize()
	if err != nil {
		return nil, err
	}
	serialized.Write(btb)

	return serialized.Bytes(), nil
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

	var baseTxBranch []Hash
	branchSize := treeDepth(uint(txCount))
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

	pb.MajorVersion = byte(majorVersion)
	pb.MinorVersion = byte(minorVersion)
	pb.Timestamp = timestamp
	pb.Nonce = nonce
	pb.Prev = prev

	pb.TransactionsCount = uint16(txCount)
	pb.BaseTransactionBranch = baseTxBranch
	pb.BaseTransaction = baseTx

	return nil
}

func GenerateGenesisBlock(network *config.Network) (*Block, error) {
	var genesisBlock Block

	genesisTransactionBytes, err := hex.DecodeString(network.GenesisCoinbaseTxHex)
	reader := bytes.NewReader(genesisTransactionBytes)

	if err != nil {
		return nil, err
	}

	if err := genesisBlock.Transaction.Deserialize(reader); err != nil {
		return nil, err
	}

	genesisBlock.MajorVersion = config.BlockMajorVersion1
	genesisBlock.MinorVersion = config.BlockMinorVersion0
	genesisBlock.Timestamp = 0
	genesisBlock.Nonce = 70

	return &genesisBlock, nil
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
func treeDepth(count uint) int {
	depth := 0

	for i := unsafe.Sizeof(count) << 2; i > 0; i >>= 1 {
		if count >> 1 > 0 {
			count >>= i
			depth += 1
		}
	}

	return depth
}
