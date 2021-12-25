package cryptonote

import "errors"

var (
	ErrAddBlockAlreadyExists              = errors.New("block already exists")
	ErrAddBlockTransactionCountNotMatch   = errors.New("transaction sizes not match")
	ErrAddBlockTransactionSizeMax         = errors.New("transaction size bigger than allowed")
	ErrAddBlockTransactionCoinbaseMaxSize = errors.New("coinbase transaction size bigger than allowed")
	ErrAddBlockTransactionDeserialization = errors.New("transaction deserialization failed")
	ErrAddBlockRejectedAsOrphaned         = errors.New("rejected as orphaned")

	ErrAddBlockFailedGetDifficulty = errors.New("failed to get difficulty for next block")

	// ErrAddBlockUnexpectedError any unpredictable error
	ErrAddBlockUnexpectedError = errors.New("unexpected error")
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
	ErrBlockValidationBlockSignatureMismatch      = errors.New("block signature mismatch")
	ErrBlockValidationCheckpointBlockHashMismatch = errors.New("checkout block hash mismatch")
	ErrBlockValidationProofOfWorkTooWeak          = errors.New("proof of work too weak")
	ErrBlockValidationTransactionAbsentInPool     = errors.New("transaction absent in pool")
	ErrBlockValidationBaseTransactionExtraMMTag   = errors.New("base transaction extra MM tag")
	ErrBlockValidationTransactionInconsistency    = errors.New("transaction inconsistency")
	ErrBlockValidationDuplicateTransaction        = errors.New("duplicate transaction")
)

var (
	ErrTransactionEmptyInputs                          = errors.New("transaction has no inputs")
	ErrTransactionInputUnknownType                     = errors.New("transaction has input with unknown type")
	ErrTransactionInputEmptyOutputUsage                = errors.New("transaction's input uses empty output")
	ErrTransactionInputInvalidDomainKeyImages          = errors.New("transaction uses key image not in the valid domain")
	ErrTransactionInputIdenticalKeyImages              = errors.New("transaction has identical key images")
	ErrTransactionInputIdenticalOutputIndexes          = errors.New("transaction has identical output indexes")
	ErrTransactionInputKeyImageAlreadySpent            = errors.New("transaction uses spent key image")
	ErrTransactionInputMultisignatureAlreadySpent      = errors.New("transaction uses spent multisignature")
	ErrTransactionInputInvalidGlobalIndex              = errors.New("transaction has input with invalid global index")
	ErrTransactionInputSpendLockedOut                  = errors.New("transaction uses locked input")
	ErrTransactionInputInvalidSignatures               = errors.New("transaction has input with invalid signature")
	ErrTransactionInputWrongSignaturesCount            = errors.New("transaction has input with wrong signatures count")
	ErrTransactionInputsAmountOverflow                 = errors.New("transaction's inputs sum overflow")
	ErrTransactionInputWrongCount                      = errors.New("wrong input count")
	ErrTransactionInputUnexpectedType                  = errors.New("wrong input type")
	ErrTransactionBaseInputWrongBlockIndex             = errors.New("base input has wrong block index")
	ErrTransactionOutputZeroAmount                     = errors.New("transaction has zero output amount")
	ErrTransactionOutputInvalidKey                     = errors.New("transaction has output with invalid key")
	ErrTransactionOutputInvalidRequiredSignaturesCount = errors.New("transaction has output with invalid signatures count")
	ErrTransactionOutputInvalidMultisignatureKey       = errors.New("transaction has output with invalid multisignature key")
	ErrTransactionOutputUnknownType                    = errors.New("transaction has unknown output type")
	ErrTransactionOutputsAmountOverflow                = errors.New("transaction has outputs amount overflow")
	ErrTransactionWrongAmount                          = errors.New("transaction wrong amount")
	ErrTransactionWrongUnlockTime                      = errors.New("transaction has wrong unlock time")
	ErrTransactionInvalidMixin                         = errors.New("transaction has wrong mixin")
	ErrTransactionExtraTooLarge                        = errors.New("transaction extra is too large")
	ErrTransactionBaseInvalidSignaturesCount           = errors.New("coinbase transactions must not have input signatures")
	ErrTransactionInputInvalidSignaturesCount          = errors.New("the number of input signatures is not correct")
	ErrTransactionOutputInvalidDecomposedAmount        = errors.New("invalid decomposed output amount (unmixable output)")
	ErrTransactionInvalidFee                           = errors.New("fee is too small and it's not a fusion transaction")
	ErrTransactionSizeTooLarge                         = errors.New("transaction is too large (in bytes)")
	ErrTransactionOutputsInvalidCount                  = errors.New("only 1 output in coinbase transaction allowed")
	ErrTransactionBaseOutputWrongType                  = errors.New("coinbase transaction can have only output key output type")
)
