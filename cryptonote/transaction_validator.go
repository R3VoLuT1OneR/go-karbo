package cryptonote

import (
	"github.com/r3volut1oner/go-karbo/config"
	log "github.com/sirupsen/logrus"
)

type TransactionValidator struct {
	// transaction that will be validated
	transaction *Transaction

	// blockHeight of the transaction
	blockHeight uint32

	// network of the transaction
	network *config.Network

	// logger entry with already transaction fields configured
	logger *log.Entry
}

func (validator *TransactionValidator) Validate() error {
	if err := validator.validateSize(); err != nil {
		return err
	}

	if err := validator.validateInputs(); err != nil {
		return err
	}

	return nil
}

func (validator *TransactionValidator) validateSize() error {
	size := validator.transaction.Size()
	maxSize := validator.network.MaxTransactionSize(validator.blockHeight)

	if size > maxSize {
		err := ErrTransactionSizeTooLarge
		validator.logger.WithFields(log.Fields{
			"transaction_size":     size,
			"transaction_maz_size": maxSize,
		}).Error(err)
		return err
	}

	return nil
}

func (validator *TransactionValidator) validateInputs() error {
	inputs := validator.transaction.Inputs

	if len(inputs) == 0 {
		err := ErrTransactionEmptyInputs
		validator.logger.Error(err)
		return err
	}

	//sumOfInputs := uint64(0)
	//keyImageSet := map[KeyImage]bool{}
	//
	//for i, input := range inputs {
	//	logger := validator.logger.WithFields(log.Fields{
	//		"transaction_input_index": i,
	//	})
	//
	//	amount := uint64(0)
	//
	//	switch input.(type) {
	//	case InputKey:
	//		input := input.(InputKey)
	//		logger := logger.WithField("transaction_input_type", "InputKey")
	//
	//		amount := input.Amount
	//
	//		if _, ok := keyImageSet[input.KeyImage]; ok {
	//			err := ErrTransactionInputIdenticalKeyImages
	//			logger.Error(err)
	//			return err
	//		}
	//
	//		keyImageSet[input.KeyImage] = true
	//
	//		if len(input.OutputIndexes) == 0 {
	//			err := ErrTransactionInputEmptyOutputUsage
	//			logger.Error(err)
	//			return err
	//		}
	//
	//		// outputIndexes are packed here, first is absolute, others are offsets to previous,
	//		// so first can be zero, others can't.
	//		// Fix discovered by Monero Lab and suggested by "fluffypony" (bitcointalk.org).
	//		// Skip this expensive validation in checkpoints zone.
	//
	//	case InputMultiSignature:
	//	default:
	//		err := ErrTransactionInputsAmountOverflow
	//		logger.Error(err)
	//		return err
	//	}
	//}

	return nil
}
