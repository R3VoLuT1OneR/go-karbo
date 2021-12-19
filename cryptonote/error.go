package cryptonote

import "errors"

var (
	ErrAddBlockAlreadyExists              = errors.New("block already exists")
	ErrAddBlockTransactionCountNotMatch   = errors.New("transaction sizes not match")
	ErrAddBlockTransactionSizeMax         = errors.New("transaction size bigger than allowed")
	ErrAddBlockTransactionCoinbaseMaxSize = errors.New("coinbase transaction size bigger than allowed")
	ErrAddBlockTransactionDeserialization = errors.New("transaction deserialization failed")
	ErrAddBlockRejectedAsOrphaned         = errors.New("rejected as orphaned")
)

var (
	ErrBlockValidationCumulativeSizeTooBig        = errors.New("cumulative size too big")
	ErrBlockValidationWrongVersion                = errors.New("wrong block version")
	ErrBlockValidationParentBlockSizeTooBig       = errors.New("parent block size too big")
	ErrBlockValidationParentBlockWrongVersion     = errors.New("parent block wrong version")
	ErrBlockValidationTimestampTooFarInFuture     = errors.New("timestamp too far in future")
	ErrBlockValidationTimestampTooFarInPast       = errors.New("timestamp too far in past")
	ErrBlockValidationDifficultyOverhead          = errors.New("difficulty overhead")
	ErrBlockValidationBlockRewardMismatch         = errors.New("block reward mismatch")
	ErrBlockValidationCheckpointBlockHashMismatch = errors.New("checkout block hash mismatch")
	ErrBlockValidationProofOfWorkTooWeak          = errors.New("proof of work too weak")
	ErrBlockValidationTransactionAbsentInPool     = errors.New("transaction absent in pool")
	ErrBlockValidationBaseTransactionExtraMMTag   = errors.New("base transaction extra MM tag")
	ErrBlockValidationTransactionInconsistency    = errors.New("transaction inconsistency")
	ErrBlockValidationDuplicateTransaction        = errors.New("duplicate transaction")
)
