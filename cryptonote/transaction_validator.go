package cryptonote

import (
	"errors"
	"github.com/r3volut1oner/go-karbo/config"
	"github.com/r3volut1oner/go-karbo/crypto"
	"github.com/r3volut1oner/go-karbo/utils"
	log "github.com/sirupsen/logrus"
	"math"
)

type blockTransactionsValidator struct {
	// bc is the link to blockchain
	bc *BlockChain

	// blockIndex of the transaction
	blockIndex uint32

	// spentKeyImages used for sharing spent key images for block validation
	spentKeyImages map[crypto.KeyImage]bool

	// spentMultisignatureGlobalIndexes used for sharing spent key images for block validation
	spentMultisignatureGlobalIndexes map[MultisigAmountIndexPair]bool

	// cumulativeFee used for sharing spent key images for block validation
	cumulativeFee uint64

	// logger entry with already transaction fields configured
	logger *log.Logger
}

type MultisigAmountIndexPair struct {
	Amount      uint64
	OutputIndex uint32
}

var LImage = crypto.KeyImage(crypto.L)
var IImage = crypto.KeyImage(crypto.I)

func NewBlockTransactionsValidator(bc *BlockChain, blockIndex uint32, logger *log.Logger) *blockTransactionsValidator {
	return &blockTransactionsValidator{
		bc:                               bc,
		blockIndex:                       blockIndex,
		spentKeyImages:                   map[crypto.KeyImage]bool{},
		spentMultisignatureGlobalIndexes: map[MultisigAmountIndexPair]bool{},
		cumulativeFee:                    uint64(0),
		logger:                           logger,
	}
}

func (validator *blockTransactionsValidator) validate(transaction *Transaction) error {
	if err := validator.validateSize(transaction); err != nil {
		return err
	}

	sumOfInputs, err := validator.validateInputs(transaction)
	if err != nil {
		return err
	}

	sumOfOutputs, err := validator.validateOutputs(transaction)
	if err != nil {
		return err
	}

	fee, isFusion, err := validator.validateTransactionFee(transaction, sumOfInputs, sumOfOutputs)
	if err != nil {
		return err
	}

	if err := validator.validateTransactionExtra(transaction, fee, isFusion); err != nil {
		return err
	}

	if err := validator.validateTransactionMixin(transaction); err != nil {
		return err
	}

	// Verify key images are not spent, ring signatures are valid, etc. We
	// do this separately from the transaction input verification, because
	// these checks are much slower to perform, so we want to fail fast on the
	// cheaper checks first.
	if err := validator.validateTransactionInputExpensive(transaction); err != nil {
		return err
	}

	return nil
}

func (validator *blockTransactionsValidator) validateSize(transaction *Transaction) error {
	size := transaction.Size()
	maxSize := validator.bc.Network.MaxTransactionSize(validator.blockIndex)
	logger := validator.logger.WithFields(log.Fields{
		"transaction_size":     size,
		"transaction_maz_size": maxSize,
	})

	if size > maxSize {
		err := ErrTransactionSizeTooLarge
		logger.Error(err)
		return err
	}

	return nil
}

func (validator *blockTransactionsValidator) validateInputs(transaction *Transaction) (uint64, error) {
	inputs := transaction.Inputs

	if len(inputs) == 0 {
		err := ErrTransactionEmptyInputs
		validator.logger.Error(err)
		return 0, err
	}

	sumOfInputs := uint64(0)
	keyImageSet := map[crypto.KeyImage]bool{}

	for i, input := range inputs {
		logger := validator.logger.WithFields(log.Fields{
			"transaction_input_index": i,
		})

		amount := uint64(0)

		switch input.(type) {
		case InputKey:
			input := input.(InputKey)
			logger := logger.WithFields(log.Fields{
				"transaction_input_type": "InputKey",
			})

			amount = input.Amount

			if _, ok := keyImageSet[input.KeyImage]; ok {
				err := ErrTransactionInputIdenticalKeyImages
				logger.Error(err)
				return 0, err
			}

			keyImageSet[input.KeyImage] = true

			if len(input.OutputIndexes) == 0 {
				err := ErrTransactionInputEmptyOutputUsage
				logger.Error(err)
				return 0, err
			}

			// outputIndexes are packed here, first is absolute, others are offsets to previous,
			// so first can be zero, others can't.
			// Fix discovered by Monero Lab and suggested by "fluffypony" (bitcointalk.org).
			// Skip this expensive validation in checkpoints zone.
			if !validator.bc.Network.Checkpoints.IsInCheckpointZone(validator.blockIndex + 1) {
				multiplied, err := input.KeyImage.ScalarMult(&LImage)

				if err != nil || *multiplied != IImage {
					err := ErrTransactionInputInvalidDomainKeyImages
					logger.Error(err)
					return 0, err
				}
			}

			if !utils.SliceIsUniqueUint32(&input.OutputIndexes) {
				err := ErrTransactionInputIdenticalOutputIndexes
				logger.Error(err)
				return 0, err
			}

			if _, ok := validator.spentKeyImages[input.KeyImage]; ok {
				err := ErrTransactionInputKeyImageAlreadySpent
				logger.Error(err)
				return 0, err
			}

			validator.spentKeyImages[input.KeyImage] = true
		case InputMultiSignature:
			input := input.(InputMultiSignature)
			logger := logger.WithFields(log.Fields{
				"transaction_input_type": "InputMultiSignature",
			})

			MOPair := MultisigAmountIndexPair{input.Amount, input.OutputIndex}
			if _, ok := validator.spentMultisignatureGlobalIndexes[MOPair]; ok {
				err := ErrTransactionInputMultisignatureAlreadySpent
				logger.Error(err)
				return 0, err
			}

			validator.spentMultisignatureGlobalIndexes[MOPair] = true

			output, unlockTime, exists :=
				validator.bc.IsMultiSignatureOutputExists(input.Amount, input.OutputIndex, validator.blockIndex)

			if !exists {
				err := ErrTransactionInputInvalidGlobalIndex
				logger.Error(err)
				return 0, err
			}

			if validator.bc.IsMultiSignatureSpent(input.Amount, input.OutputIndex, validator.blockIndex) {
				err := ErrTransactionInputMultisignatureAlreadySpent
				logger.Error(err)
				return 0, err
			}

			if !validator.bc.IsTransactionSpendTimeUnlocked(unlockTime, validator.blockIndex) {
				err := ErrTransactionInputSpendLockedOut
				logger.Error(err)
				return 0, err
			}

			if output.RequiredSignaturesCount != input.SignatureCount {
				err := ErrTransactionInputWrongSignaturesCount
				logger.Error(err)
				return 0, err
			}
		default:
			err := ErrTransactionInputUnknownType
			logger.Error(err)
			return 0, err
		}

		if math.MaxUint64-amount < sumOfInputs {
			err := ErrTransactionInputsAmountOverflow
			logger.Error(err)
			return 0, err
		}

		sumOfInputs += amount
	}

	return sumOfInputs, nil
}

func (validator *blockTransactionsValidator) validateOutputs(transaction *Transaction) (uint64, error) {
	sumOfOutputs := uint64(0)

	for i, output := range transaction.Outputs {
		logger := validator.logger.WithFields(log.Fields{
			"transaction_output_index": i,
		})

		if output.Amount == 0 {
			err := ErrTransactionOutputZeroAmount
			logger.Error(err)
			return 0, err
		}

		if validator.blockIndex >= config.UpgradeHeightV5 {
			if !validator.bc.Network.IsValidDecomposedAmount(output.Amount) {
				err := ErrTransactionOutputInvalidDecomposedAmount
				logger.Error(err)
				return 0, err
			}
		}

		switch output.Target.(type) {
		case OutputKey:
			outputTarget := output.Target.(OutputKey)

			if !outputTarget.Check() {
				err := ErrTransactionOutputInvalidKey
				logger.Error(err)
				return 0, err
			}
		case OutputMultisignature:
			outputTarget := output.Target.(OutputMultisignature)

			if outputTarget.RequiredSignaturesCount > byte(len(outputTarget.Keys)) {
				err := ErrTransactionOutputInvalidRequiredSignaturesCount
				logger.Error(err)
				return 0, err
			}

			for _, key := range outputTarget.Keys {
				logger := logger.WithFields(log.Fields{
					"transaction_output_key_index": i,
				})

				if !key.Check() {
					err := ErrTransactionOutputInvalidMultisignatureKey
					logger.Error(err)
					return 0, err
				}
			}
		default:
			err := ErrTransactionOutputUnknownType
			logger.Error(err)
			return 0, err
		}

		if math.MaxUint64-output.Amount < sumOfOutputs {
			err := ErrTransactionOutputsAmountOverflow
			logger.Error(err)
			return 0, err
		}

		sumOfOutputs += output.Amount
	}

	return sumOfOutputs, nil
}

// Pre-requisite - Call validateTransactionInputs() and validateTransactionOutputs()
// to ensure m_sumOfInputs and m_sumOfOutputs is set
func (validator *blockTransactionsValidator) validateTransactionFee(transaction *Transaction, sumOfInputs, sumOfOutputs uint64) (fee uint64, isFusion bool, err error) {
	logger := validator.logger.WithFields(log.Fields{
		"sum_of_inputs":  sumOfInputs,
		"sum_of_outputs": sumOfOutputs,
	})

	if sumOfInputs == 0 || sumOfOutputs == 0 {
		err = errors.New("sum of inputs or outputs are zero")
		logger.Error(err)
		return
	}

	if sumOfOutputs > sumOfInputs {
		err = ErrTransactionWrongAmount
		logger.Error(err)
		return
	}

	fee = sumOfInputs - sumOfOutputs

	isFusion = fee == 0 && validator.IsFusionTransaction(transaction)

	if !isFusion {
		h := validator.blockIndex
		minFee := validator.bc.Network.MinimalFeeValidator(h)

		if fee < minFee {
			logger := validator.logger.WithFields(log.Fields{
				"transaction_fee":     fee,
				"transaction_min_fee": minFee,
			})

			err = ErrTransactionInvalidFee
			logger.Error(err)
			return
		}
	}

	validator.cumulativeFee += fee

	return
}

func (validator *blockTransactionsValidator) validateTransactionExtra(transaction *Transaction, fee uint64, isFusion bool) error {
	minFee := validator.bc.Network.MinimalFee(validator.blockIndex)
	extraSize := uint64(len(transaction.Extra))
	feePerByte := validator.bc.Network.GetFeePerByte(extraSize, minFee)
	minFeeWithExtra := uint64(0)
	min := minFee + feePerByte

	if validator.blockIndex > config.UpgradeHeightV4s2 && validator.blockIndex < config.UpgradeHeightV4s3 {
		minFeeWithExtra = min - ((min * 20) / 100)
	} else if validator.blockIndex >= config.UpgradeHeightV4s3 {
		minFeeWithExtra = min
	}

	if minFeeWithExtra != 0 && !isFusion && fee < minFeeWithExtra {
		logger := validator.logger.WithFields(log.Fields{
			"transaction_extra_size":    extraSize,
			"transaction_extra_min_fee": minFeeWithExtra,
			"transaction_fee":           fee,
		})

		err := ErrTransactionInvalidFee
		logger.Error(err)
		return err
	}

	return nil
}

func (validator *blockTransactionsValidator) validateTransactionMixin(transaction *Transaction) error {
	mixin := countMixin(transaction)
	minMixin := validator.bc.Network.MinMixin()
	maxMixin := validator.bc.Network.MaxMixin()

	if (validator.blockIndex > config.UpgradeHeightV3s1 && mixin > maxMixin) ||
		(validator.blockIndex > config.UpgradeHeightV4 && mixin < minMixin && mixin != 1) {
		logger := validator.logger.WithFields(log.Fields{
			"transaction_mixin":     mixin,
			"transaction_max_mixin": maxMixin,
			"transaction_min_mixin": minMixin,
		})

		err := ErrTransactionInvalidMixin
		logger.Error(err)
		return err
	}

	return nil
}

func (validator *blockTransactionsValidator) validateTransactionInputExpensive(transaction *Transaction) error {
	// Don't need to do expensive transaction validation for transactions
	// in a checkpoints range - they are assumed valid, and the transaction
	// hash would change thus invalidation the checkpoints if not.

	if validator.bc.Network.Checkpoints.IsInCheckpointZone(validator.blockIndex) {
		return nil
	}

	prefixHash := transaction.TransactionPrefix.Hash()

	for inputIndex, input := range transaction.Inputs {
		logger := validator.logger.WithFields(log.Fields{
			"transaction_input_index": inputIndex,
		})

		switch input.(type) {
		case InputKey:
			input := input.(InputKey)
			logger := validator.logger.WithFields(log.Fields{
				"transaction_input_type": "InputKey",
			})

			if validator.bc.IsSpent(input.KeyImage, validator.blockIndex) {
				err := ErrTransactionInputKeyImageAlreadySpent
				logger.Error(err)
				return err
			}

			globalIndexes := make([]uint32, len(input.OutputIndexes))
			globalIndexes[0] = input.OutputIndexes[0]
			for i := 1; i < len(input.OutputIndexes); i++ {
				globalIndexes[i] = globalIndexes[i-1] + input.OutputIndexes[i]
			}

			outputKeys, err := validator.bc.ExtractKeyOutputKeys(input.Amount, validator.blockIndex, globalIndexes)

			// Handle extract key output keys error
			if err != nil {
				switch err {
				case ErrExtractOutputKeyInvalidGlobalIndex:
					err = ErrTransactionInputInvalidGlobalIndex
				case ErrExtractOutputKeyLocked:
					err = ErrTransactionInputSpendLockedOut
				default:
					err = ErrTransactionUnknownError
				}

				logger.Error(err)
				return err
			}

			sigs := transaction.TransactionSignatures[inputIndex]

			if len(outputKeys) != len(sigs) {
				err := ErrTransactionInputInvalidSignaturesCount
				logger.Error(err)
				return err
			}

			if !crypto.CheckRingSignature(prefixHash, &input.KeyImage, &outputKeys, &sigs, config.KeyImageCheckingBlockIndex) {
				err := ErrTransactionInputInvalidSignatures
				logger.Error(err)
				return err
			}

		case InputMultiSignature:
			// input := input.(InputMultiSignature)
			logger := validator.logger.WithFields(log.Fields{
				"transaction_input_type": "InputMultiSignature",
			})

			err := ErrTransactionMultiSignaturesNotImplemented
			logger.Error(err)
			return err
		default:
			err := ErrTransactionInputUnknownType
			logger.Error(err)
			return err
		}
	}

	return nil
}

func (validator *blockTransactionsValidator) IsFusionTransaction(transaction *Transaction) bool {
	inputsAmounts := getInputsAmounts(transaction)
	outputsAmounts := getOutputsAmounts(transaction)

	if validator.blockIndex <= config.UpgradeHeightV3 {
		size := transaction.Size()
		maxSize := validator.bc.Network.FusionMaxTxSize(validator.blockIndex)

		if size > maxSize {
			logger := validator.logger.WithFields(log.Fields{
				"fusion_transaction_size":     size,
				"fusion_transaction_max_size": maxSize,
			})

			err := ErrTransactionFusionMaxSize
			logger.Error(err)
			return false
		}
	}

	minInputCount := int(validator.bc.Network.FusionTxMinInputCount())
	if len(inputsAmounts) < minInputCount {
		logger := validator.logger.WithFields(log.Fields{
			"fusion_transaction_input_len": len(inputsAmounts),
			"fusion_transaction_input_min": minInputCount,
		})

		err := ErrTransactionFusionInputCountsLessThanMinimum
		logger.Error(err)
		return false
	}

	minRatio := int(validator.bc.Network.FusionTxMinInOutCountRatio())
	if len(inputsAmounts) < len(outputsAmounts)*minRatio {
		logger := validator.logger.WithFields(log.Fields{
			"fusion_transaction_output_len": len(outputsAmounts),
			"fusion_transaction_input_len":  len(inputsAmounts),
			"fusion_transaction_min_ratio":  minRatio,
		})

		err := ErrTransactionFusionRatioInvalid
		logger.Error(err)
		return false
	}

	inputAmount := uint64(0)
	for i, amount := range inputsAmounts {
		if validator.blockIndex < config.UpgradeHeightV4 {
			dustThreshold := validator.bc.Network.DefaultDustThreshold()

			if amount < dustThreshold {
				logger := validator.logger.WithFields(log.Fields{
					"dust_threshold":     dustThreshold,
					"input_amount_index": i,
					"input_amount":       amount,
				})

				err := ErrTransactionFusionAmountLessThenThreshold
				logger.Error(err)
				return false
			}
		}

		inputAmount += amount
	}

	dustThreshold := uint64(0)
	if validator.blockIndex < config.UpgradeHeightV4 {
		dustThreshold = validator.bc.Network.DefaultDustThreshold()
	}

	expectedOutputsAmount := validator.bc.DecomposeAmount(inputAmount, dustThreshold)
	// Why do we need sorting here?
	//sort.Slice(expectedOutputsAmount, func(i, j int) bool {
	//	return expectedOutputsAmount[i] < expectedOutputsAmount[j]
	//})

	if len(expectedOutputsAmount) != len(outputsAmounts) {
		logger := validator.logger.WithFields(log.Fields{
			"expected_outputs_amount_len": len(expectedOutputsAmount),
			"input_amounts_len":           len(outputsAmounts),
		})

		err := ErrTransactionFusionDecomposedNotMatch
		logger.Error(err)
		return false
	}

	return true
}

func countMixin(transaction *Transaction) int {
	mixin := 0
	for _, input := range transaction.Inputs {
		if _, ok := input.(InputKey); ok {
			continue
		}

		curMixin := len(input.(InputKey).OutputIndexes)
		if curMixin > mixin {
			mixin = curMixin
		}
	}

	return mixin
}

func getOutputsAmounts(transaction *Transaction) []uint64 {
	outputAmounts := make([]uint64, len(transaction.Outputs))

	for i, output := range transaction.Outputs {
		outputAmounts[i] = output.Amount
	}

	return outputAmounts
}

func getInputsAmounts(transaction *Transaction) []uint64 {
	inputsAmounts := make([]uint64, len(transaction.Inputs))

	for i, input := range transaction.Inputs {
		switch input.(type) {
		case InputKey:
			inputsAmounts[i] = input.(InputKey).Amount
		case InputMultiSignature:
			inputsAmounts[i] = input.(InputMultiSignature).Amount
		}
	}

	return inputsAmounts
}
