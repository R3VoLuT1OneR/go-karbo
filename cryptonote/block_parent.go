package cryptonote

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/r3volut1oner/go-karbo/crypto"
)

type ParentBlock struct {
	MajorVersion byte
	MinorVersion byte
	Timestamp    uint64
	Nonce        uint32
	Prev         crypto.Hash

	TransactionsCount     uint16
	BaseTransactionBranch crypto.HashList
	BaseTransaction       BaseTransaction
	BlockchainBranch      []crypto.Hash
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
	var prev crypto.Hash
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

	var baseTxBranch crypto.HashList
	branchSize := treeDepth(int(txCount))
	for i := 0; i < branchSize; i++ {
		var th crypto.Hash
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
		return fmt.Errorf("failed to parse extra: %w", err)
	}

	if tef.MiningTag == nil {
		return errors.New("can't get extra merge mining tag")
	}

	if tef.MiningTag.Depth > 8*32 {
		return errors.New("wrong merge mining tag depth")
	}

	var blockchainBranch crypto.HashList
	for i := uint64(0); i < tef.MiningTag.Depth; i++ {
		var h crypto.Hash
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
