package karbowanecd

import (
	"github.com/ybbus/jsonrpc"
)

type karbowanecd struct {
	jsonrpc jsonrpc.RPCClient
}

// RPCKarbowanecd client
type RPCKarbowanecd interface {
	// GetBlockCount returns the current chain height.
	GetBlockCount() (count int, err error)

	// GetBlockHash returns block hash by its height.
	GetBlockHash(height int) (hash string, err error)

	// GetBlockTemplate returns blocktemplate with an empty "hole" for nonce.
	// Size of reserve and address can be specified
	GetBlockTemplate(size int, addr string) (template BlockTemplate, err error)

	// GetBlockHeaderByHash block header by given block hash.
	GetBlockHeaderByHash(hash string) (header BlockHeader, err error)

	// GetBlockHeaderByHeight returns block header by given block height
	GetBlockHeaderByHeight(height int) (header BlockHeader, err error)

	// GetBlockTimestamp returns the timestamp of a block.
	GetBlockTimestamp(height int) (timestamp int, err error)

	// GetBlockByHeight returns information on a single block by height.
	GetBlockByHeight(height int) (block Block, err error)

	// GetBlockByHash returns information on a single block by its hash.
	GetBlockByHash(hash string) (block Block, err error)

	// GetBlocksByHeights returns blocks by list of heights.
	GetBlocksByHeights(heights []int) (blocks []Block, err error)

	// GetBlocksByHashes returns blocks by list of hashes.
	GetBlocksByHashes(hashes []string) (blocks []Block, err error)

	// GetBlocksHashesByTimestamps returns block by timestamps
	GetBlocksHashesByTimestamps(begin, end, limit int) (hashes []string, count int, err error)

	// GetBlocksList returns a list of blocks short info, starting from height and how much
	GetBlocksList(height int, count int) (shortBlocks []ShortBlock, err error)

	// GetAltBlocksList returns a list of alt. blocks short info.
	GetAltBlocksList() (altBlocks []ShortBlock, err error)

	// GetLastBlockHeader returns the last block header.
	GetLastBlockHeader() (header BlockHeader, err error)

	// GetTransaction returns information on a single transaction by its hash.
	GetTransaction(hash string) (transaction Transaction, err error)

	// GetTransactionsPool the list of short details of transactions present in mempool.
	GetTransactionsPool() (transactions []MempoolTransaction, err error)

	// gettransactionsbypaymentid() returns an array of short data of trandactions containing given Payment ID.
	GetTransactionsByPaymentID(paymentID string) (transactions []ShortTransaction, err error)

	// GetTransactionHashesByPaymentID returns hashes of transactions containing given Payment ID.
	GetTransactionHashesByPaymentID(paymentID string) (hashes []string, err error)

	// GetTransactionByHashes returns details of the transactions found by given hashes.
	GetTransactionsByHashes(hashes []string) (transactions []Transaction, err error)

	// GetCurrencyID returns unique currency identifier.
	GetCurrencyID() (id string, err error)

	// GetStatsByHeights returns stats by heights.
	GetStatsByHeights(heights []int) (stats []BlockStats, err error)

	// GetStatsInRange returns stats in heights range.
	GetStatsInRange(start, end int) (stats []BlockStats, err error)

	// CheckTransactionKey allows to check payment by Private Transaction Key.
	CheckTransactionKey(hash, privateKey, address string) (check CheckTransaction, err error)

	// CheckTransactionByViewKey allows to check payment by recipient's View Secret Key.
	CheckTransactionByViewKey(hash, viewKey, address string) (check CheckTransactionByView, err error)

	// CheckTransactionProof allows to check payment by special Transaction Proof
	// without revealing Private Transaction Key by sender.
	CheckTransactionProof(hash, signature, destinationAddress string) (check CheckTransactionProof, err error)

	// CheckReserveProof allows to check the proof of reserve which proves that given
	// address possess/possessed the said amount at given height.
	CheckReserveProof(addr, sig, message string, height int) (check CheckReserveProof, err error)

	// ValidateAddress allows to check if given public address is valid.
	ValidateAddress(address string) (validation AddressValidation, err error)

	// VerifyMessage allows to verify message signed by wallet keys.
	VerifyMessage(message, address, signature string) (result bool, err error)

	// SubmitBlock submits mined block.
	SubmitBlock(blockBlob string) (err error)
}

// BlockTemplate represetns block template response
type BlockTemplate struct {
	BlockTemplateBlob string `json:"blocktemplate_blob"` //	Blocktemplate with empty "hole" for nonce	string
	BlockHashingBlob  string `json:"blockhashing_blob"`  // Block hashing
	Difficulty        int    `json:"difficulty"`         // Difficulty of the network
	Height            int    `json:"height"`             //	Chain height of the network
	ReservedOffset    int    `json:"reserved_offset"`    // Offset reserved
}

// Block represents blockchain block
type Block struct {
	AlreadyGeneratedCoins        int     `json:"alreadyGeneratedCoins"`        // total number of coins generated in the network upto that block
	AlreadyGeneratedTransactions int     `json:"alreadyGeneratedTransactions"` // total number of transactions present in the network upto that block
	BaseReward                   int     `json:"baseReward"`                   // calculated reward
	BlockSize                    int     `json:"blockSize"`                    // size of the block
	CumulativeDifficulty         int     `json:"cumulativeDifficulty"`         // cumulative difficulty is the sum of all block's difficulties up to this block including
	Depth                        int     `json:"depth"`                        // height away from the known top block
	Difficulty                   int     `json:"difficulty"`                   // difficulty of the requested block
	EffectiveSizeMedian          int     `json:"effectiveSizeMedian"`          // fixed constant for max size of block
	Hash                         string  `json:"hash"`                         // hash of the requested block
	Index                        int     `json:"index"`                        // index of the requested block
	IsOrphaned                   bool    `json:"isOrphaned"`                   // whether the requested block was an orphan or not
	MajorVersion                 int     `json:"majorVersion"`                 // -
	MinorVersion                 int     `json:"minorVersion"`                 // -
	Nonce                        int     `json:"nonce"`                        // -
	Penalty                      float32 `json:"penalty"`                      // penalty in block reward determined for deviation
	PrevBlockHash                string  `json:"prevBlockHash"`                // hash of the previous block
	Reward                       int     `json:"reward"`                       // total reward of the block after removing penalty
	SizeMedian                   int     `json:"sizeMedian"`                   // calculated median size from last 100 blocks
	Timestamp                    int     `json:"timestamp"`                    // the time at which the block is occured on chain since Unix epoch
	TotalFeeAmount               int     `json:"totalFeeAmount"`               // total fees for the transactions in the block
	TransactionsCumulativeSize   int     `json:"transactionsCumulativeSize"`   // total sum of size of all transactions in the block

	Transactions []*Transaction `json:"transactions"` // Array of transactions in the block
}

// ShortBlock represents short block details
type ShortBlock struct {
	Timestamp         int    `json:"timestamp"`          //	The timestamp of the block
	Height            int    `json:"height"`             //	The height of the block
	Hash              string `json:"hash"`               //	The hash of the block
	TransactionsCount int    `json:"transactions_count"` //	The number of transactions in block
	CumulativeSize    int    `json:"cumulative_size"`    //	Total size of the block
	Difficulty        int    `json:"difficulty"`         //	Block's difficulty
	MinFee            int    `json:"min_fee"`            //	Minimal transaction fee at the height of the block
}

// Transaction represetnts blockchain transaction
type Transaction struct {
	BlockHash          string                  `json:"blockHash"`          // hash of block containing this transaction
	BlockIndex         int                     `json:"blockIndex"`         // index of block containing this transaction
	Extra              TransactionExtra        `json:"extra"`              // transaction extra which can be any information	json
	Fee                int                     `json:"fee"`                // total fees of the transaction
	Hash               string                  `json:"hash"`               // hash of the transaction
	InBlockchain       bool                    `json:"inBlockchain"`       // wherer thransaction is included in block
	Inputs             []*TransactionInput     `json:"inputs"`             // inputs of the transaction json
	Mixin              int                     `json:"mixin"`              // mixin of the transaction
	Outputs            []*TransactionOutput    `json:"outputs"`            // outputs of the transaction	json
	PaymentID          string                  `json:"paymentId"`          // payment Id of the transaction
	Signatures         []*TransactionSignature `json:"signatures"`         // array of transaction signatures
	SignaturesSize     int                     `json:"signaturesSize"`     // the size of the signatures array
	Size               int                     `json:"size"`               // total size of the transaction
	Timestamp          int                     `json:"timestamp"`          // timestamp of the block that includes transaction
	TotalInputsAmount  int                     `json:"totalInputsAmount"`  // total amount on input side of the transaction
	TotalOutputsAmount int                     `json:"totalOutputsAmount"` // total amount present in the transaction
	UnlockTime         int                     `json:"unlockTime"`         // delay in unlocking the amount
	Version            int                     `json:"version"`
}

// MempoolTransaction represents transaction in mempool
type MempoolTransaction struct {
	ShortTransaction
	ReceiveTime int `json:"receive_time"` //	the timestamp when transaction was received by queried node
}

// ShortTransaction represetns short transaction details
type ShortTransaction struct {
	Hash      string `json:"hash"`       //	the hash of transaction
	Fee       int    `json:"fee"`        //	the network fee of transaction
	AmountOut int    `json:"amount_out"` //	total amount in transaction
	Size      int    `json:"size"`       //	the size of transaction
}

// BlockHeader represents header in blockchain block
type BlockHeader struct {
	BlockSize    int    `json:"block_size"`    //	size of the block
	Depth        int    `json:"depth"`         //	height away from the known top block
	Difficulty   int    `json:"difficulty"`    //	difficulty of the requested block
	Hash         string `json:"hash"`          //	hash of the requested block
	Height       int    `json:"height"`        //	height of the requested block
	MajorVersion int    `json:"major_version"` //	-
	MinorVersion int    `json:"minor_version"` //	-
	Nonce        int    `json:"nonce"`         //	-
	NumTxs       int    `json:"num_txs"`       //	Number of transactions in the block
	OrphanStatus bool   `json:"orphan_status"` //	whether the requested block was an orphan or not
	PrevHash     string `json:"prev_hash"`     //	hash of the previous block
	Reward       int    `json:"reward"`        //	reward of the block
	Timestamp    int    `json:"timestamp"`     //	the time at which the block is occured on chain since Unix epoch
}

// TransactionExtra represents extra transaction data
type TransactionExtra struct {
	Nonce     []int  `json:"nonce"`
	PublicKey string `json:"publicKey"`
	Raw       string `json:"raw"`
	Size      int    `json:"size"`
}

// TransactionInput represents transaction input
type TransactionInput struct {
	Type string `json:"type"`
	Data struct {
		// type == 'ff' only
		Amount int `json:"amount,omitempty"`

		Input struct {
			// type == '02'the
			Amount     int    `json:"amount"`
			KImage     string `json:"k_image"`
			KeyOffsets []int  `json:"key_offsets"`

			// type == 'ff'
			Heigh int `json:"height"`
		} `json:"input"`

		Mixin int `json:"mixin"`

		Outputs []struct {
			Number          int    `json:"number"`
			TransactionHash string `json:"transactionHash"`
		} `json:"outputs"`
	} `json:"data,omitempty"`
}

// TransactionOutput represents transaction output
type TransactionOutput struct {
	GlobalIndex int     `json:"globalIndex"`
	Output      *Output `json:"output"`
}

// TransactionSignature represents transaction signature
type TransactionSignature struct {
	First  int    `json:"first"`
	Second string `json:"second"`
}

// Output represents blockchain output
type Output struct {
	Amount int `json:"amount"`
	Target struct {
		Data struct {
			Key string `json:"key"`
		} `json:"data"`
		Type string `json:"type"`
	} `json:"target"`
}

// BlockStats represents block stats
type BlockStats struct {
	Height                int `json:"height"`
	AlreadyGeneratedCoins int `json:"already_generated_coins"`
	TransactionsCount     int `json:"transactions_count"`
	BlockSize             int `json:"block_size"`
	Difficulty            int `json:"difficulty"`
	Reward                int `json:"reward"`
	Timestamp             int `json:"timestamp"`
}

// CheckTransaction represents check transaction response
type CheckTransaction struct {
	Amount  int      `json:"amount"`  // the amount received by destination address
	Outputs []Output `json:"outputs"` // outputs that belong to destination address
}

// CheckTransactionByView represents response of check transaction by view
type CheckTransactionByView struct {
	Amount        int      `json:"amount"`        // the amount received by destination address
	Outputs       []Output `json:"outputs"`       // outputs that belong to destination address
	Confirmations int      `json:"confirmations"` //	the number of network confirmations
}

// CheckTransactionProof represents response of check transaction proof
type CheckTransactionProof struct {
	SignatureValid bool     `json:"signature_valid"` //	whether signature is valid
	ReceivedAmount int      `json:"received_amount"` //	the amount received by destination address
	Confirmations  int      `json:"confirmations"`   //	the number of network confirmations
	Outputs        []Output `json:"outputs"`         //	outputs that belong to destination address
}

// CheckReserveProof represents response of check reserve proof
type CheckReserveProof struct {
	Good   bool `json:"good"`   //	whether signature is valid
	Total  int  `json:"total"`  //	total amount proved
	Spent  int  `json:"spent"`  //	spent amount
	Locked int  `json:"locked"` //	locked amount
}

// AddressValidation represents address validation response
type AddressValidation struct {
	IsValid        bool   `json:"is_valid"`         //	whether given address is valid
	Address        string `json:"address"`          //	address, encoded from decoded public keys
	SpendPublicKey string `json:"spend_public_key"` //	decoded spend public key of the address
	ViewPublicKey  string `json:"view_public_key"`  //	decoded view public key of the address
}

// NewClient return walletd rpc client
func NewClient(endpoint string) RPCKarbowanecd {
	return NewClientWithOpts(endpoint, nil)
}

// NewClientWithOpts return walletd rpc client with additional options
func NewClientWithOpts(endpoint string, opts *jsonrpc.RPCClientOpts) RPCKarbowanecd {
	rpcClient := jsonrpc.NewClientWithOpts(endpoint, opts)

	client := &karbowanecd{
		jsonrpc: rpcClient,
	}

	return client
}

func (client *karbowanecd) GetBlockCount() (height int, err error) {
	var rsp struct {
		Count int `json:"count"`
	}

	if err = client.jsonrpc.CallFor(&rsp, "getblockcount", struct{}{}); err == nil {
		height = int(rsp.Count)
	}

	return
}

func (client *karbowanecd) GetBlockHash(height int) (hash string, err error) {
	err = client.jsonrpc.CallFor(&hash, "getblockhash", [1]int{height})
	// err = client.jsonrpc.CallFor(&hash, "getblockhash", struct {
	// 	Height int `json:"height"`
	// }{height})

	return
}

func (client *karbowanecd) GetBlockTemplate(size int, addr string) (template BlockTemplate, err error) {
	err = client.jsonrpc.CallFor(&template, "getblocktemplate", &struct {
		Size int    `json:"reserve_size"`   // Size of the reserve to be specified
		Addr string `json:"wallet_address"` // Valid wallet address
	}{size, addr})

	return
}

func (client *karbowanecd) GetBlockHeaderByHash(hash string) (header BlockHeader, err error) {
	rsp := struct {
		BlockHeader *BlockHeader `json:"block_header"`
	}{&header}

	req := struct {
		Hash string `json:"hash"`
	}{hash}

	err = client.jsonrpc.CallFor(&rsp, "getblockheaderbyhash", req)

	return
}

func (client *karbowanecd) GetBlockHeaderByHeight(height int) (header BlockHeader, err error) {
	rsp := struct {
		BlockHeader BlockHeader `json:"block_header"`
	}{header}

	req := struct {
		Height int `json:"height"`
	}{height}

	if err = client.jsonrpc.CallFor(&rsp, "getblockheaderbyheight", req); err == nil {
		header = rsp.BlockHeader
	}

	return
}

func (client *karbowanecd) GetBlockTimestamp(height int) (timestamp int, err error) {
	var rsp struct {
		Timestamp int    `json:"timestamp"`
		Status    string `json:"status"`
	}

	req := struct {
		Height int `json:"height"`
	}{height}

	if err = client.jsonrpc.CallFor(&rsp, "getblocktimestamp", req); err == nil {
		timestamp = rsp.Timestamp
	}

	return
}

func (client *karbowanecd) GetBlockByHeight(height int) (block Block, err error) {
	rsp := struct {
		Block Block `json:"block"`
	}{block}

	req := struct {
		Height int `json:"blockHeight"`
	}{height}

	if err = client.jsonrpc.CallFor(&rsp, "getblockbyheight", req); err == nil {
		block = rsp.Block
	}

	return
}

func (client *karbowanecd) GetBlockByHash(hash string) (block Block, err error) {
	rsp := struct {
		Block Block `json:"block"`
	}{block}

	req := struct {
		Hash string `json:"hash"`
	}{hash}

	if err = client.jsonrpc.CallFor(&rsp, "getblockbyhash", req); err == nil {
		block = rsp.Block
	}

	return
}

func (client *karbowanecd) GetBlocksByHeights(heights []int) (blocks []Block, err error) {
	rsp := struct {
		Blocks *[]Block `json:"blocks"`
	}{&blocks}

	req := struct {
		BlockHeights []int `json:"blockHeights"`
	}{heights}

	err = client.jsonrpc.CallFor(&rsp, "getblocksbyheights", req)

	return
}

func (client *karbowanecd) GetBlocksByHashes(hashes []string) (blocks []Block, err error) {
	rsp := struct {
		Blocks *[]Block `json:"blocks"`
	}{&blocks}

	req := struct {
		Hashes []string `json:"blockHashes"`
	}{hashes}

	err = client.jsonrpc.CallFor(&rsp, "getblocksbyhashes", req)

	return
}

func (client *karbowanecd) GetBlocksHashesByTimestamps(begin, end, limit int) (hashes []string, count int, err error) {
	rsp := struct {
		Hashes *[]string `json:"blockHashes"`
		Count  *int      `json:"count"`
	}{&hashes, &count}

	req := struct {
		TimestampBegin int `json:"timestampBegin"`
		TimestampEnd   int `json:"timestampEnd"`
		Limit          int `json:"limit"`
	}{begin, end, limit}

	err = client.jsonrpc.CallFor(&rsp, "getblockshashesbytimestamps", req)

	return
}

func (client *karbowanecd) GetBlocksList(height int, count int) (shortBlocks []ShortBlock, err error) {
	// 10 is the default blocks limit
	if count == 0 {
		count = 10
	}

	rsp := struct {
		ShortBlocks *[]ShortBlock `json:"blocks"`
	}{&shortBlocks}

	req := struct {
		Height int `json:"height"`
		Count  int `json:"count"`
	}{height, count}

	err = client.jsonrpc.CallFor(&rsp, "getblockslist", req)

	return
}

func (client *karbowanecd) GetAltBlocksList() (altBlocks []ShortBlock, err error) {
	rsp := struct {
		AltBlocks *[]ShortBlock `json:"alt_blocks"`
	}{&altBlocks}

	err = client.jsonrpc.CallFor(&rsp, "getaltblockslist")

	return
}

func (client *karbowanecd) GetLastBlockHeader() (header BlockHeader, err error) {
	rsp := struct {
		BlockHeader *BlockHeader `json:"block_header"`
	}{&header}

	err = client.jsonrpc.CallFor(&rsp, "getlastblockheader")

	return
}

func (client *karbowanecd) GetTransaction(hash string) (transaction Transaction, err error) {
	rsp := struct {
		Transaction *Transaction `json:"transaction"`
	}{&transaction}

	req := struct {
		Hash string `json:"hash"`
	}{hash}

	err = client.jsonrpc.CallFor(&rsp, "gettransaction", req)

	return
}

func (client *karbowanecd) GetTransactionsPool() (transactions []MempoolTransaction, err error) {
	rsp := struct {
		Transactions *[]MempoolTransaction `json:"transactions"`
	}{&transactions}

	err = client.jsonrpc.CallFor(&rsp, "gettransactionspool")

	return
}

func (client *karbowanecd) GetTransactionsByPaymentID(paymentID string) (transactions []ShortTransaction, err error) {
	rsp := struct {
		Transactions *[]ShortTransaction `json:"transactions"`
	}{&transactions}

	req := struct {
		PaymentID string `json:"payment_id"`
	}{paymentID}

	err = client.jsonrpc.CallFor(&rsp, "gettransactionsbypaymentid", req)

	return
}

func (client *karbowanecd) GetTransactionHashesByPaymentID(paymentID string) (hashes []string, err error) {
	rsp := struct {
		TransactionHashes *[]string `json:"transactionHashes"`
	}{&hashes}

	req := struct {
		PaymentID string `json:"paymentId"`
	}{paymentID}

	err = client.jsonrpc.CallFor(&rsp, "gettransactionhashesbypaymentid", req)

	return
}

func (client *karbowanecd) GetTransactionsByHashes(hashes []string) (transactions []Transaction, err error) {
	rsp := struct {
		Transactions *[]Transaction `json:"transactions"`
	}{&transactions}

	req := struct {
		Hashes []string `json:"transactionHashes"`
	}{hashes}

	err = client.jsonrpc.CallFor(&rsp, "gettransactionsbyhashes", req)

	return
}

func (client *karbowanecd) GetCurrencyID() (id string, err error) {
	rsp := struct {
		CurrencyIDBlob *string `json:"currency_id_blob"`
	}{&id}

	err = client.jsonrpc.CallFor(&rsp, "getcurrencyid")

	return
}

func (client *karbowanecd) GetStatsByHeights(heights []int) (stats []BlockStats, err error) {
	rsp := struct {
		Stats *[]BlockStats `json:"stats"`
	}{&stats}

	req := struct {
		Height []int `json:"heights"`
	}{heights}

	err = client.jsonrpc.CallFor(&rsp, "getstatsbyheights", req)

	return
}

func (client *karbowanecd) GetStatsInRange(start, end int) (stats []BlockStats, err error) {
	rsp := struct {
		Stats *[]BlockStats `json:"stats"`
	}{&stats}

	req := struct {
		StartHeight int `json:"start_height"`
		EndHeight   int `json:"end_height"`
	}{start, end}

	err = client.jsonrpc.CallFor(&rsp, "getstatsinrange", req)

	return
}

func (client *karbowanecd) CheckTransactionKey(hash, privateKey, address string) (check CheckTransaction, err error) {
	req := struct {
		Hash       string `json:"transaction_id"`
		PrivateKey string `json:"transaction_key"`
		Address    string `json:"address"`
	}{hash, privateKey, address}

	err = client.jsonrpc.CallFor(&check, "checktransactionkey", req)

	return
}

func (client *karbowanecd) CheckTransactionByViewKey(hash, viewKey, address string) (check CheckTransactionByView, err error) {
	req := struct {
		Hash       string `json:"transaction_id"`
		PrivateKey string `json:"view_key"`
		Address    string `json:"address"`
	}{hash, viewKey, address}

	err = client.jsonrpc.CallFor(&check, "checktransactionbyviewkey", req)

	return
}

func (client *karbowanecd) CheckTransactionProof(hash, signature, destinationAddress string) (check CheckTransactionProof, err error) {
	req := struct {
		Hash               string `json:"transaction_id"`
		Signature          string `json:"signature"`
		DestinationAddress string `json:"destination_address"`
	}{hash, signature, destinationAddress}

	err = client.jsonrpc.CallFor(&check, "checktransactionproof", req)

	return
}

func (client *karbowanecd) CheckReserveProof(addr, sig, message string, height int) (check CheckReserveProof, err error) {
	req := struct {
		Address   string `json:"address"`
		Signature string `json:"signature"`
		Message   string `json:"message,omitempty"`
		Height    int    `json:"height,omitempty"`
	}{addr, sig, message, height}

	err = client.jsonrpc.CallFor(&check, "checkreserveproof", req)

	return
}

func (client *karbowanecd) ValidateAddress(address string) (validation AddressValidation, err error) {
	req := struct {
		Address string `json:"address"`
	}{address}

	err = client.jsonrpc.CallFor(&validation, "validateaddress", req)

	return
}

func (client *karbowanecd) VerifyMessage(message, address, signature string) (result bool, err error) {
	rsp := struct {
		SigValid *bool `json:"sig_valid"`
	}{&result}

	req := struct {
		Message   string `json:"message"`
		Address   string `json:"address"`
		Signature string `json:"signature"`
	}{message, address, signature}

	err = client.jsonrpc.CallFor(&rsp, "verifymessage", req)

	return
}

func (client *karbowanecd) SubmitBlock(blockBlob string) (err error) {
	req := struct {
		BlockBlob string `json:"block_blob"`
	}{blockBlob}

	_, err = client.jsonrpc.Call("submitblock", req)

	return
}
