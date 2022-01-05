package cryptonote

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/crypto"
	"unsafe"
)

// Block consists of three parts:
// - block header
// - base transaction body
// - list of transaction identifiers (hashes)
//
// The list starts with the number of transaction identifiers that it
// contains.
type Block struct {
	// Each block starts with a block header.
	BlockHeader

	// Each valid block contains a single base transaction. The base
	// transaction's validity depends on the block height due to the
	// following reasons:
	//   - the emission rule is generally defined as a function of time;
	//   - without the block height field, two base transactions could
	//     be indistinguishable as they can have the same hash (see [BH] for
	//     a description of a similar problem in Bitcoin).
	BaseTransaction Transaction

	// A transaction identifier is a transaction body hashed with the Keccak
	// hash function. The list starts with the number of identifiers and is
	// followed by the identifiers themselves if it is not empty.
	TransactionsHashes []crypto.Hash

	// ParentBlock was introduced in the 2 and 3, block versions and removed and newest.
	// Pointer is used to save space in blocks with version that is not using it.
	Parent *ParentBlock

	// Signature introduced in the 5th block version.
	// It is the proof that blocked was mined by the owner of the rewarded address.
	// Pointer is used to save space in blocks with version that is not using it.
	Signature *crypto.Signature

	// Next variables are used only for caching in runtime
	hash             *crypto.Hash
	hashTransactions *crypto.Hash
}

// BlockHeader represents the metadata of the block
//
// The major version defines the block header parsing rules (i.e. block header format) and is incremented with
// each block header format update.
// The minor version defines the interpretation details that are not related to block header parsing.
type BlockHeader struct {
	// MajorVersion of the block
	// Blockchain involves so the block may change the look
	MajorVersion byte

	// MinorVersion of the block
	MinorVersion byte

	// Timestamp of when the block was mined
	Timestamp uint64

	// PreviousBlockHash have a unique identifier of the previous block
	// Used for building a chain as it is play a role of parent block
	PreviousBlockHash crypto.Hash

	Nonce uint32
}

// Hash returns hash of the block.
// It is used as unique identifier of the block.
func (b *Block) Hash() *crypto.Hash {
	if b.hash == nil {
		var allBytesBuffer bytes.Buffer

		/**
		 * First part of the hashing bytes
		 */
		allBytesBuffer.Write(b.HashingBytes())

		if b.MajorVersion == config.BlockMajorVersion2 || b.MajorVersion == config.BlockMajorVersion3 {
			allBytesBuffer.Write(b.Parent.serialize(true))
		}

		/**
		 * Create final hash bytes, by appending hash bytes length and the hash bytes
		 */
		allBytes := allBytesBuffer.Bytes()
		allBytesCount := make([]byte, binary.MaxVarintLen64)
		written := binary.PutUvarint(allBytesCount, uint64(len(allBytes)))

		var h bytes.Buffer
		h.Write(allBytesCount[:written])
		h.Write(allBytes)

		b.hash = new(crypto.Hash)
		b.hash.FromBytes(h.Bytes())
	}

	return b.hash
}

func (b *Block) Index() uint32 {
	if len(b.BaseTransaction.Inputs) == 1 {
		i := b.BaseTransaction.Inputs[0]

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

	if majorVersion >= config.BlockMajorVersion5 {
		sig := &crypto.Signature{}
		if err := sig.Deserialize(r); err != nil {
			return err
		}

		b.Signature = sig
	}

	if majorVersion == config.BlockMajorVersion2 || majorVersion == config.BlockMajorVersion3 {
		parentBlock := &ParentBlock{}
		if err := parentBlock.deserialize(r); err != nil {
			return err
		}

		b.Parent = parentBlock
	}

	if err := b.BaseTransaction.Deserialize(r); err != nil {
		return err
	}

	hashesCount, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}

	hl := crypto.HashList{}
	for i := uint64(0); i < hashesCount; i++ {
		var h crypto.Hash
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

	majorVersion := b.BlockHeader.MajorVersion

	if majorVersion >= config.BlockMajorVersion5 {
		sSignature, _ := b.Signature.Serialize()
		serialized.Write(sSignature)
	}

	if majorVersion == config.BlockMajorVersion2 || majorVersion == config.BlockMajorVersion3 {
		serialized.Write(b.Parent.serialize(false))
	}

	serialized.Write(b.BaseTransaction.Serialize())

	buf := make([]byte, binary.MaxVarintLen64)
	written := binary.PutUvarint(buf, uint64(len(b.TransactionsHashes)))
	serialized.Write(buf[:written])

	for i := 0; i < len(b.TransactionsHashes); i++ {
		serialized.Write(b.TransactionsHashes[i][:])
	}

	return serialized.Bytes()
}

// HashingBytes is a copy of the C++ implementation of getBlockHashingBinaryArray method.
// it is used for first step of the serialization
func (b *Block) HashingBytes() []byte {
	var allBytesBuffer bytes.Buffer

	/**
	 * Write block header bytes
	 */
	allBytesBuffer.Write(b.BlockHeader.serialize())

	/**
	 * Write merkle root hash bytes
	 */
	baseTransactionHash := b.BaseTransaction.Hash()
	hl := crypto.HashList{*baseTransactionHash}
	hl = append(hl, b.TransactionsHashes...)
	allBytesBuffer.Write(hl.MerkleRootHash()[:])

	/**
	 * Write transactions number
	 */
	transactionCount := make([]byte, binary.MaxVarintLen64)
	written := binary.PutUvarint(transactionCount, uint64(len(hl)))
	allBytesBuffer.Write(transactionCount[:written])

	return allBytesBuffer.Bytes()
}

func (h *BlockHeader) deserialize(reader *bytes.Reader) error {
	var prev crypto.Hash
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
