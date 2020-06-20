package walletd

import (
	"errors"

	"github.com/ybbus/jsonrpc"
)

// ResBalance is response when we ret
type ResBalance struct {
	AvailableBalance int64
	LockedAmount     int64
}

// ReqReset used for Reset request
type ReqReset struct {
	ViewSecretKey string `json:"viewSecretKey,omitempty"`
	ScanHeight    int    `json:"scanHeight,omitempty"`
}

// ReqCreateAddress used for request CreateAddress
//
// SecretSpendKey - If was specified, RPC Wallet creates spend address
// PublicSpendKey - If was specified, RPC Wallet creates view address
// ScanHeight - Height from which start synchronization
// Reset - Determines reset wallet or not, after creation. By default do reset.
type ReqCreateAddress struct {
	SecretSpendKey string `json:"secretSpendKey,omitempty"`
	PublicSpendKey string `json:"publicSpendKey,omitempty"`
	ScanHeight     int    `json:"scanHeight,omitempty"`
	Reset          bool   `json:"reset,omitempty"`
}

// ResCreateAddress used for response CreateAddress
//
// Address - new base58 encoded address
type ResCreateAddress struct {
	Address string `json:"address"`
}

// ReqCreateAddressList used for request CreateAddressList
//
// SpendSecretKeys - Required. Array of private spend keys.
// ScanHeights - Array of height, determines a height to start wallet sync.
// Reset - Array of reset values needed to be reset. Do reset by default
type ReqCreateAddressList struct {
	SpendSecretKeys []string `json:"spendSecretKeys"`
	ScanHeights     []int    `json:"scanHeights,omitempty"`
	Reset           []bool   `json:"reset,omitempty"`
}

// ResCreateAddressList represents response
//
// Addresses - Array of strings, where each string is a created address
type ResCreateAddressList struct {
	Addresses []string `json:"addresses"`
}

// ReqExport used for Export request
type ReqExport struct {
	FileName string `json:"fileName"`
}

// ResSpendKeys used for spend keys response
type ResSpendKeys struct {
	SpendSecretKey string `json:"spendSecretKey"`
	SpendPublicKey string `json:"spendPublicKey"`
}

// ReqGetBlockHashes used as request for get block hashes
// BlockCount - Number of blocks to process
// FirstBlockIndex - Starting height
type ReqGetBlockHashes struct {
	BlockCount      int `json:"blockCount"`
	FirstBlockIndex int `json:"firstBlockIndex"`
}

// ResGetBlockHashes list of block hashes
type ResGetBlockHashes struct {
	BlockHashes []string `json:"blockHashes"`
}

// ReqGetTransactionHashes used as request
//
// Addresses - (Optional) Array of strings, where each string is an address.
// BlockHash - Hash of the starting block.
// FirstBlockIndex - Starting height.
// BlockCount - Number of blocks to return transaction hashes from.
// PaymentId - (Optional) Valid payment_id.
type ReqGetTransactionHashes struct {
	Addresses       []string `json:"addresses,omitempty"`
	BlockHash       string   `json:"blockHash,omitempty"`
	FirstBlockIndex int      `json:"firstBlockIndex,omitempty"`
	BlockCount      int      `json:"blockCount"`
	PaymentID       string   `json:"paymentId,omitempty"`
}

// ResGetTransactionHashes used as response
// Array that contains: blockHash - string - hash of the block which contains
// transaction hashes; transactionHashes - array - array of strings, where
// each string is a transaction hash
type ResGetTransactionHashes struct {
	Items []struct {
		TransactionHashes []string `json:"transactionHashes"`
		BlockHash         string   `json:"blockHash"`
	} `json:"items"`
}

// ReqGetTransactions req for transaction
type ReqGetTransactions struct {
	Addresses       []string `json:"addresses,omitempty"`
	BlockHash       string   `json:"blockHash,omitempty"`
	FirstBlockIndex int      `json:"firstBlockIndex,omitempty"`
	BlockCount      int      `json:"blockCount"`
	PaymentID       string   `json:"paymentId,omitempty"`
}

// ResGetTransactions res for transaction
type ResGetTransactions struct {
	Items []struct {
		Transacations []TransactionResponse `json:"transactions"`
		BlockHash     string                `json:"blockHash"`
	} `json:"items"`
}

// ReqGetTransactionProof request
type ReqGetTransactionProof struct {
	TransactionHash      string `json:"transactionHash"`
	DestinationAddress   string `json:"destinationAddress"`
	TransactionSecretKey string `json:"transactionSecretKey"`
}

// ResGetTransactionProof response
type ResGetTransactionProof struct {
	TransactionProof string `json:"transactionProof"`
}

// ReqSendTransaction request
type ReqSendTransaction struct {
	// Array of strings, where each string is an address to take the funds from
	Addresses []string `json:"addresses,omitempty"`

	// Array that contains: address - string; amount - int64
	Transfers []struct {
		// Address
		Address string `json:"address"`
		// Amount
		Amount uint64 `json:"amount"`
	} `json:"transfers"`

	// Transaction fee. Minimal fee in Karbowanec network is .0001 KRB.
	// This parameter should be specified in minimal available KRB units.
	// For example, if your fee is .01 KRB, you should pass it as 10000000000
	Fee int `json:"fee"`

	// Height of the block until which transaction is going to be locked for
	// spending.
	UnlockTime uint64 `json:"unlockTime,omitempty"`

	// Privacy level (a discrete number from 1 to infinity). Level 6 and higher
	// is recommended.
	Anonymity int `json:"anonymity"`

	// String of variable length. Can contain A-Z, 0-9 characters.
	Extra string `json:"extra,omitempty"`

	// Payment ID
	PaymentID string `json:"paymentId,omitempty"`

	// Valid and existing in this container address.
	ChangeAddress string `json:"changeAddress,omitempty"`
}

// ResSendTransaction response
type ResSendTransaction struct {
	TransactionHash string `json:"transactionHash"`
}

// ReqCreateDelayedTransaction request
type ReqCreateDelayedTransaction struct {
	// Array of strings, where each string is an address to take the funds from
	Addresses []string `json:"addresses,omitempty"`

	// Array that contains: address - string; amount - int64
	Transfers []struct {
		// Address
		Address string `json:"address"`
		// Amount
		Amount uint64 `json:"amount"`
	} `json:"transfers"`

	// Transaction fee. Minimal fee in Karbowanec network is .0001 KRB.
	// This parameter should be specified in minimal available KRB units.
	// For example, if your fee is .01 KRB, you should pass it as 10000000000
	Fee int `json:"fee"`

	// Height of the block until which transaction is going to be locked for
	// spending.
	UnlockTime uint64 `json:"unlockTime,omitempty"`

	// Privacy level (a discrete number from 1 to infinity). Level 6 and higher
	// is recommended.
	Anonymity int `json:"anonymity"`

	// String of variable length. Can contain A-Z, 0-9 characters.
	Extra string `json:"extra,omitempty"`

	// Payment ID
	PaymentID string `json:"paymentId,omitempty"`

	// Valid and existing in this container address.
	ChangeAddress string `json:"changeAddress,omitempty"`
}

// ResCreateDelayedTransaction response
type ResCreateDelayedTransaction struct {
	TransactionHash string `json:"transactionHash"`
}

// TransactionResponse represents transaction in response
type TransactionResponse struct {
	// "fee":1000000,
	Fee int `json:"fee"`

	// "extra":"0130b4472974f2deb9fae7d8fd6602b26396379f3fa05cca2430e10e9e60179f42",
	Extra string `json:"extra"`

	// "timestamp":0,
	Timestamp int `json:"timestamp"`

	// "blockIndex":4294967295,
	BlockIndex int `json:"blockIndex"`

	// "state":0,
	State int `json:"state"`

	// "transactionHash":"92423b0857d36bd172b3f2effbd47ea477bfe0618a50c29d475542c6d5d1b835",
	TransactionHash string `json:"transactionHash"`

	// "amount":-1703701,
	Amount int `json:"amount"`

	// "unlockTime":0,
	UnlockTime int `json:"unlockTime"`

	Transfers []struct {

		// "amount":123456,
		Amount int `json:"amount"`

		// "type":0,
		Type int `json:"type"`

		// "address":"KiQxu9U3F7vdGggu4NQ3CKDhk59vMQyMaFbLtu7TU4TdUkNtuJufqpo67r2e5j5p44SBsBBygaRdmeB4gwH9CF1C3zufGWd"
		Address string `json:"address"`
	} `json:"transfers"`

	// "paymentId":"",
	PaymentID string `json:"paymentId"`

	// "isBase":False
	IsBase bool `json:"isBase"`
}

// ReqGetUnconfirmedTransactionHashes request
type ReqGetUnconfirmedTransactionHashes struct {
	Addresses []string `json:"addresses"`
}

// ResGetUnconfirmedTransactionHashes response
type ResGetUnconfirmedTransactionHashes struct {
	TransactionHashes []string `json:"transactionHashes"`
}

// ResGetTransaction response
type ResGetTransaction struct {
	Transaction TransactionResponse `json:"transaction"`
}

// ResGetTransactionSecretKey response
type ResGetTransactionSecretKey struct {
	TransactionSecretKey string `json:"transactionSecretKey"`
}

// ResGetDelayedTransactionHashes response
type ResGetDelayedTransactionHashes struct {
	TransactionHashes []string `json:"transactionHashes"`
}

// ResGetViewKey response
type ResGetViewKey struct {
	ViewSecretKey string `json:"viewSecretKey"`
}

// ResGetMnemonicSeed response
type ResGetMnemonicSeed struct {
	MnemonicSeed string `json:"mnemonicSeed"`
}

// ResGetStatus response
type ResGetStatus struct {
	// Node's known number of blocks
	BlockCount uint32 `json:"blockCount"`

	//Maximum known number of blocks of all seeds that are connected to the node
	KnownBlockCount uint32 `json:"knownBlockCount"`

	// Local (synchronized) blocks count
	LocalDaemonBlockCount uint32 `json:"localDaemonBlockCount"`

	// Hash of the last known block
	LastBlockHash string `json:"lastBlockHash"`

	// Connected peers number
	PeerCount uint32 `json:"peerCount"`

	// Current minimum transaction fee in atomic units. Do not use received value
	// 'as is', but to round it up to one of two first digits after leading zeroes
	// or double it to make sure tx will pass in case of minimalFee fluctuations.
	MinimalFee uint64 `json:"minimalFee"`

	// The version of the Payment Gate software
	Version string `json:"version"`
}

// ResGetAddresses response
type ResGetAddresses struct {
	Addresses []string `json:"addresses"`
}

// ResGetAddressesCount response
type ResGetAddressesCount struct {
	AddressesCount uint32 `json:"addressesCount"`
}

// ReqSendFusionTransaction request
type ReqSendFusionTransaction struct {
	// Value that determines which outputs will be optimized. Only the outputs,
	// lesser than the threshold value, will be included into a fusion transaction.
	Threshold uint64 `json:"threshold"`

	// Privacy level (a discrete number from 1 to infinity). Level 6 and higher
	// is recommended.
	Anonymity uint64 `json:"anonymity"`

	// Array of strings, where each string is an address to take the funds from.
	Addresses []string `json:"addresses,omitempty"`

	// An address that the optimized funds will be sent to. Valid and existing
	// in this container address.
	DestinationAddress string `json:"destinationAddress,omitempty"`
}

// ResSendFusionTransaction response
type ResSendFusionTransaction struct {
	TransactionHash string `json:"transactionHash"`
}

// ReqEstimateFusion request
type ReqEstimateFusion struct {
	// Value that determines which outputs will be optimized. Only the outputs,
	// lesser than the threshold value, will be included into a fusion transaction.
	Threshold uint64 `json:"threshold"`

	// Array of strings, where each string is an address to take the funds from.
	Addresses []string `json:"addresses"`
}

// ResEstimateFusion response
type ResEstimateFusion struct {
	// Total number of unspent outputs of the specified addresses.
	TotalOutputCount int `json:"totalOutputCount"`

	// Number of outputs that can be optimized.
	FusionReadyCount int `json:"fusionReadyCount"`
}

// ResValidateAddress response
type ResValidateAddress struct {
	Address        string `json:"address"`
	IsValid        bool   `json:"isValid"`
	SpendPublicKey string `json:"spendPublicKey"`
	ViewPublicKey  string `json:"viewPublicKey"`
}

// ReqGetReserveProof request
type ReqGetReserveProof struct {
	Address string `json:"address"`
	Message string `json:"message"`
	Amount  uint64 `json:"amount"`
}

// ResGetReserveProof response
type ResGetReserveProof struct {
	ReserveProof string `json:"reserveProof"`
}

// ReqSignMessage request
type ReqSignMessage struct {
	// Public address for keys used to sign message
	Address string `json:"address,omitempty"`

	// The message to sign
	Message string `json:"message"`
}

// ResSignMessage response
type ResSignMessage struct {
	// Address, used to sing the message (useful in case it was omitted in request)
	Adddress string `json:"address"`

	// The signature generated with wallet keys
	Signature string `json:"signature"`
}

// ReqVerifyMessage request
type ReqVerifyMessage struct {
	// Public address for keys used to sign message
	Address string `json:"address"`

	// The message to sign
	Message string `json:"message"`

	// The signature
	Signature string `json:"signature"`
}

// ResVerifyMessage response
type ResVerifyMessage struct {
	// Denotes whether signature of message is valid or not
	IsValid bool `json:"isValid"`
}

// RPCWalletd sends requests to Walletd RPC server
type RPCWalletd interface {

	// Save current state of wallet
	Save() error

	// Reset allows to re-sync wallet
	//
	// Important: If the view_secret_key was not pointed out reset() methods
	// resets the wallet and re-syncs it. If the view_secret_key argument
	// was pointed out reset() method substitutes the existing wallet with
	// a new one with a specified view_secret_key and creates an address for it.
	Reset(opts *ReqReset) error

	// Export wallet to container file
	Export(opts *ReqExport) error

	// CreateAddress creates an additional address in wallet
	//
	// If opts == nil new spend key generated
	//
	// Provide spend key to generate address for the keys. One of the parameters
	// must be provided.
	CreateAddress(opts *ReqCreateAddress) (*ResCreateAddress, error)

	// CreateAddressList method creates additional addresses in your wallet from
	// the provided list of private keys.
	CreateAddressList(opts *ReqCreateAddressList) (*ResCreateAddressList, error)

	// DeleteAddress removes address from the container
	DeleteAddress(address string) error

	// GetSpendKeys returns address spend keys
	GetSpendKeys(address string) (*ResSpendKeys, error)

	// GetBalance returns balance for specific address
	// If address == nil provided will provide default address balance
	GetBalance(address string) (*ResBalance, error)

	// GetBlockHashes() method returns an array of block hashes for a specified
	// block range.
	GetBlockHashes(opts *ReqGetBlockHashes) (*ResGetBlockHashes, error)

	// GetTransactionHashes method returns an array of block and transaction hashes.
	// Transaction consists of transfers. Transfer is an amount-address pair.
	// There could be several transfers in a single transaction.
	//
	// Note: if paymentId parameter is set, getTransactionHashes() method returns
	// transaction hashes of transactions that contain specified payment_id.
	// (in the set block range)
	//
	// Note: if addresses parameter is set, getTransactionHashes() method returns
	// transaction hashes of transactions that contain transfer from at least one
	// of specified addresses.
	//
	// Note: if both above mentioned parameters are set, getTransactionHashes()
	// method returns transaction hashes of transactions that contain both
	// specified payment_id and transfer from at least one of specified addresses.
	GetTransactionHashes(opts *ReqGetTransactionHashes) (*ResGetTransactionHashes, error)

	// GetTransactions method returns an array of block and transaction hashes.
	// Transaction consists of transfers. Transfer is an amount-address pair.
	// There could be several transfers in a single transaction.
	//
	// Note: if paymentId parameter is set, getTransactions() method returns transactions
	// that contain specified payment_id. (in the set block range)
	//
	// Note: if addresses parameter is set, getTransactions() method returns transactions
	// that contain transfer from at least one of specified addresses.
	//
	// Note: if both above mentioned parameters are set, getTransactions() method returns
	// transactions that contain both specified payment_id and transfer from at least one
	// of specified addresses.
	GetTransactions(opts *ReqGetTransactions) (*ResGetTransactions, error)

	// GetUnconfirmedTransactionHashes method returns information about the current
	// unconfirmed transaction pool or for a specified addresses.
	// Transaction consists of transfers. Transfer is an amount-address pair.
	// There could be several transfers in a single transaction.
	// Note: if addresses parameter is set, getUnconfirmedTransactionHashes() method
	// returns transactions that contain transfer from at least one of specified addresses.
	GetUnconfirmedTransactionHashes(opts *ReqGetUnconfirmedTransactionHashes) (*ResGetUnconfirmedTransactionHashes, error)

	// GetTransaction method returns information about a particular transaction.
	// Transaction consists of transfers. Transfer is an amount-address pair.
	// There could be several transfers in a single transaction.
	GetTransaction(transaction string) (*ResGetTransaction, error)

	// GetTransactionSecretKey get secret application key
	GetTransactionSecretKey(transaction string) (*ResGetTransactionSecretKey, error)

	// GetTransactionProof get transaction proof
	GetTransactionProof(opts *ReqGetTransactionProof) (*ResGetTransactionProof, error)

	// SendTransaction sendTransaction() method allows you to send transaction to one or
	// several addresses. Also, it allows you to use a payment_id for a transaction to a
	// single address.
	//
	// Note: if container contains only 1 address, changeAddress field can be left empty
	// and the change is going to be sent to this address
	//
	// Note: if addresses field contains only 1 address, changeAddress can be left empty
	// and the change is going to be sent to this address
	//
	// Note: in the rest of the cases, changeAddress field is mandatory and must contain
	// an address.
	SendTransaction(opts *ReqSendTransaction) (*ResSendTransaction, error)

	// CreateDelayedTransaction method creates a delayed transaction.
	// Such transactions are not sent into the network automatically and should be pushed
	// using SendDelayedTransaction method.
	//
	// Note: if container contains only 1 address, changeAddress field can be left empty
	// and the change is going to be sent to this address
	//
	// Note: if addresses field contains only 1 address, changeAddress can be left empty
	// and the change is going to be sent to this address
	//
	// Note: in the rest of the cases, changeAddress field is mandatory and must contain
	// an address.
	//
	// Note: outputs that were used for this transactions will be locked until the
	// transaction is sent or cancelled
	CreateDelayedTransaction(opts *ReqCreateDelayedTransaction) (*ResCreateDelayedTransaction, error)

	// GetDelayedTransactionHashes method returns hashes of delayed transactions.
	GetDelayedTransactionHashes() (*ResGetDelayedTransactionHashes, error)

	// DeleteDelayedTransaction method deletes a specified delayed transaction.
	DeleteDelayedTransaction(transaction string) error

	// SendDelayedTransaction method sends a specified delayed transaction.
	SendDelayedTransaction(transaction string) error

	// GetViewKey method returns your view key.
	GetViewKey() (*ResGetViewKey, error)

	// GetMnemonicSeed method returns your first address' spend secret key and
	// derived from this spend key common private view key of the container if
	// it was created with option --deterministic.
	GetMnemonicSeed(address string) (*ResGetMnemonicSeed, error)

	// GetStatus method returns information about the current RPC Wallet state:
	// block_count, known_block_count, last_block_hash, peer_count an minimal_fee.
	GetStatus() (*ResGetStatus, error)

	// GetAddresses method returns an array of your RPC Wallet's addresses.
	GetAddresses() (*ResGetAddresses, error)

	// ResGetAddressesCount method returns count of addresses
	GetAddressesCount() (*ResGetAddressesCount, error)

	// SendFusionTransaction method allows you to send a fusion transaction, by taking
	// funds from selected addresses and transferring them to the destination address.
	//
	// If there aren't any outputs that can be optimized, SendFusionTransaction will
	// return an error. You can use estimateFusion to check the outputs, available for
	// the optimization.
	//
	// Note: if container contains only 1 address, destinationAddress field can be
	// left empty and the funds are going to be sent to this address.
	//
	// Note: if addresses field contains only 1 address, destinationAddress can be
	// left empty and the funds are going to be sent to this address.
	//
	// Note: in the rest of the cases, destinationAddress field is mandatory and must
	// contain an address.
	SendFusionTransaction(opts *ReqSendFusionTransaction) (*ResSendFusionTransaction, error)

	// EstimateFusion method counts the number of unspent outputs of the specified
	// addresses and returns how many of those outputs can be optimized.
	//
	// This method is used to understand if a fusion transaction can be created.
	// If fusionReadyCount returns a value = 0, then a fusion transaction cannot be created.
	EstimateFusion(*ReqEstimateFusion) (*ResEstimateFusion, error)

	// ValidateAddress method allows you to check if provided address is valid.
	ValidateAddress(address string) (*ResValidateAddress, error)

	// GetReserveProof get reserve proof
	GetReserveProof(opts *ReqGetReserveProof) (*ResGetReserveProof, error)

	// signMessage() allows to sign message with wallet keys.
	SignMessage(opts *ReqSignMessage) (*ResSignMessage, error)

	// VerifyMessage method is used to verify signed message.
	VerifyMessage(opts *ReqVerifyMessage) (*ResVerifyMessage, error)
}

// NewClient return walletd rpc client
func NewClient(endpoint string) RPCWalletd {
	return NewClientWithOpts(endpoint, nil)
}

// NewClientWithOpts return walletd rpc client with additional options
func NewClientWithOpts(endpoint string, opts *jsonrpc.RPCClientOpts) RPCWalletd {
	rpcClient := jsonrpc.NewClientWithOpts(endpoint, opts)

	client := &walletClient{
		rpcClient: rpcClient,
	}

	return client
}

type walletClient struct {
	rpcClient jsonrpc.RPCClient
}

func (client *walletClient) Save() error {
	_, err := client.rpcClient.Call("save")

	return err
}

func (client *walletClient) Reset(opts *ReqReset) error {
	var err error

	if opts == nil {
		_, err = client.rpcClient.Call("reset")
	} else {
		_, err = client.rpcClient.Call("reset", opts)
	}

	return err
}

func (client *walletClient) Export(opts *ReqExport) error {
	_, err := client.rpcClient.Call("export", opts)

	return err
}

func (client *walletClient) CreateAddress(opts *ReqCreateAddress) (*ResCreateAddress, error) {
	var response *ResCreateAddress
	var err error

	if opts == nil {
		err = client.rpcClient.CallFor(&response, "createAddress")
	} else {
		err = client.rpcClient.CallFor(&response, "createAddress", opts)
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) CreateAddressList(opts *ReqCreateAddressList) (*ResCreateAddressList, error) {
	var response *ResCreateAddressList

	err := client.rpcClient.CallFor(&response, "createAddressList", opts)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) DeleteAddress(address string) error {
	_, err := client.rpcClient.Call("deleteAddress", &map[string]string{
		"address": address,
	})

	return err
}

func (client *walletClient) GetSpendKeys(address string) (*ResSpendKeys, error) {
	var response *ResSpendKeys

	err := client.rpcClient.CallFor(&response, "getSpendKeys", &map[string]string{
		"address": address,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) GetBalance(address string) (*ResBalance, error) {
	var response *ResBalance
	var err error

	if address == "" {
		err = client.rpcClient.CallFor(&response, "getBalance")
	} else {
		err = client.rpcClient.CallFor(&response, "getBalance", &map[string]string{
			"address": address,
		})
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) GetBlockHashes(opts *ReqGetBlockHashes) (*ResGetBlockHashes, error) {
	var response *ResGetBlockHashes

	if err := client.rpcClient.CallFor(&response, "getBlockHashes", opts); err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) GetTransactionHashes(opts *ReqGetTransactionHashes) (*ResGetTransactionHashes, error) {
	var response *ResGetTransactionHashes

	if opts.BlockHash != "" && opts.FirstBlockIndex != 0 {
		return nil, errors.New("one of 'BlockHash' or 'FirstBlockIndex' must be provided")
	}

	if err := client.rpcClient.CallFor(&response, "getTransactionHashes", opts); err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) GetTransactions(opts *ReqGetTransactions) (*ResGetTransactions, error) {
	var response *ResGetTransactions

	if opts.BlockHash != "" && opts.FirstBlockIndex != 0 {
		return nil, errors.New("one of 'BlockHash' or 'FirstBlockIndex' must be provided")
	}

	if err := client.rpcClient.CallFor(&response, "getTransactions", opts); err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) GetUnconfirmedTransactionHashes(opts *ReqGetUnconfirmedTransactionHashes) (*ResGetUnconfirmedTransactionHashes, error) {
	var response *ResGetUnconfirmedTransactionHashes

	if err := client.rpcClient.CallFor(&response, "getUnconfirmedTransactionHashes", opts); err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) GetTransaction(transaction string) (*ResGetTransaction, error) {
	var response *ResGetTransaction

	err := client.rpcClient.CallFor(&response, "getTransaction", &map[string]string{
		"transactionHash": transaction,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) GetTransactionSecretKey(transaction string) (*ResGetTransactionSecretKey, error) {
	var response *ResGetTransactionSecretKey

	err := client.rpcClient.CallFor(&response, "getTransactionSecretKey", &map[string]string{
		"transaction": transaction,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) GetTransactionProof(opts *ReqGetTransactionProof) (*ResGetTransactionProof, error) {
	var response *ResGetTransactionProof

	if err := client.rpcClient.CallFor(&response, "getTransactionProof", opts); err != nil {
		return response, nil
	}

	return response, nil
}

func (client *walletClient) SendTransaction(opts *ReqSendTransaction) (*ResSendTransaction, error) {
	var response *ResSendTransaction

	if err := client.rpcClient.CallFor(&response, "sendTransaction", opts); err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) CreateDelayedTransaction(opts *ReqCreateDelayedTransaction) (*ResCreateDelayedTransaction, error) {
	var response *ResCreateDelayedTransaction

	if err := client.rpcClient.CallFor(&response, "createDelayedTransaction", opts); err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) GetDelayedTransactionHashes() (*ResGetDelayedTransactionHashes, error) {
	var response *ResGetDelayedTransactionHashes

	if err := client.rpcClient.CallFor(&response, "getDelayedTransactionHashes"); err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) DeleteDelayedTransaction(transaction string) error {
	_, err := client.rpcClient.Call("deleteDelayedTransaction", map[string]string{
		"transactionHash": transaction,
	})

	return err
}

func (client *walletClient) SendDelayedTransaction(transaction string) error {
	_, err := client.rpcClient.Call("sendDelayedTransaction", map[string]string{
		"transactionHash": transaction,
	})

	return err
}

func (client *walletClient) GetViewKey() (*ResGetViewKey, error) {
	var response *ResGetViewKey

	if err := client.rpcClient.CallFor(&response, "getViewKey"); err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) GetMnemonicSeed(address string) (*ResGetMnemonicSeed, error) {
	var response *ResGetMnemonicSeed

	err := client.rpcClient.CallFor(&response, "getMnemonicSeed", map[string]string{
		"address": address,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) GetStatus() (response *ResGetStatus, err error) {
	if err = client.rpcClient.CallFor(&response, "getStatus"); err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) GetAddresses() (response *ResGetAddresses, err error) {
	if err = client.rpcClient.CallFor(&response, "getAddresses"); err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) GetAddressesCount() (response *ResGetAddressesCount, err error) {
	if err = client.rpcClient.CallFor(&response, "getAddressesCount"); err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) SendFusionTransaction(opts *ReqSendFusionTransaction) (response *ResSendFusionTransaction, err error) {
	if err = client.rpcClient.CallFor(&response, "sendFusionTransaction", opts); err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) EstimateFusion(opts *ReqEstimateFusion) (response *ResEstimateFusion, err error) {
	if err = client.rpcClient.CallFor(&response, "estimateFusion", opts); err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) ValidateAddress(address string) (response *ResValidateAddress, err error) {
	err = client.rpcClient.CallFor(&response, "validateAddress", map[string]string{
		"address": address,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) GetReserveProof(opts *ReqGetReserveProof) (response *ResGetReserveProof, err error) {
	if err = client.rpcClient.CallFor(&response, "reserveProof", opts); err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) SignMessage(opts *ReqSignMessage) (response *ResSignMessage, err error) {
	if err = client.rpcClient.CallFor(&response, "signMessage", opts); err != nil {
		return nil, err
	}

	return response, nil
}

func (client *walletClient) VerifyMessage(opts *ReqVerifyMessage) (response *ResVerifyMessage, err error) {
	if err = client.rpcClient.CallFor(&response, "verifyMessage", opts); err != nil {
		return nil, err
	}

	return response, nil
}
