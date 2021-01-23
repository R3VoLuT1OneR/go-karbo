package cryptonote

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/r3volut1oner/go-karbo/config"
)

type ParentBlock struct {
	MajorVersion byte
	MinorVersion byte
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

type BlockTemplate struct {
	BlockHeader

	Parent 				ParentBlock
	Transaction 		Transaction
	TransactionsHashes 	[]*Hash

	hash *Hash
	hashTransactions *Hash
}

func (b *BlockTemplate) Hash() (*Hash, error) {
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

		hl := hashList{baseTransactionHash}
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

func GenerateGenesisBlock(network *config.Network) (*BlockTemplate, error) {
	var genesisBlock BlockTemplate

	genesisTransactionBytes, err := hex.DecodeString(network.GenesisCoinbaseTxHex)
	if err != nil {
		return nil, err
	}

	if err := genesisBlock.Transaction.Deserialize(&genesisTransactionBytes); err != nil {
		return nil, err
	}

	genesisBlock.MajorVersion = config.BlockMajorVersion1
	genesisBlock.MinorVersion = config.BlockMinorVersion0
	genesisBlock.Timestamp = 0
	genesisBlock.Nonce = 70

	return &genesisBlock, nil
}
