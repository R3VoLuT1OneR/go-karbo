package karbowanecd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ybbus/jsonrpc"
)

var testAddress = "KfXkT5VmdqmA7bWqSH37p87hSXBdTpTogN4mGHPARUSJaLse6jbXaVbVkLs3DwcmuD88xfu835Zvh6qBPCUXw6CHK8koDCt"

var testBlockHeader = map[string]interface{}{
	"depth":         1,
	"difficulty":    65198,
	"hash":          "9a8be83...",
	"height":        123456,
	"major_version": 1,
	"minor_version": 0,
	"nonce":         2358499061,
	"orphan_status": false,
	"prev_hash":     "dde56b7e...",
	"reward":        44090506423186,
	"timestamp":     1356589561,
}

var rawBlock = `{
	"alreadyGeneratedCoins":7822291504630269000,
	"alreadyGeneratedTransactions":842837,
	"baseReward":8307330332565,
	"blockSize":23716,
	"cumulativeDifficulty":2277860049927920,
	"depth":6,
	"difficulty":14377384679,
	"effectiveSizeMedian":1000000,
	"hash":"4cc94394440e64461766d0e7c675498da934d70a63b143606246f0e6e26d1bde",
	"index":399666,
	"isOrphaned":false,
	"majorVersion":4,
	"minorVersion":0,
	"nonce":3653353692,
	"penalty":0,
	"prevBlockHash":"72199371e92ee3dd95f0993f49a82d3dc7b7175dc004ad9c4437167a7138bd13",
	"reward":8507330332565,
	"sizeMedian":532,
	"timestamp":1568118204,
	"totalFeeAmount":200000000000,
	"transactions":[` + rawTransaction + `,` + rawTransaction + `],
	"transactionsCumulativeSize":23640
}`

var rawShortBlock = `{
	"cumulative_size":405,
	"difficulty":18484080,
	"hash":"f9940a120ca47d23e014078d7bbfd59f6ce4f20831a4e011b02ff73fcf0372da",
	"height":100000,
	"min_fee":100000000000,
	"timestamp":1492672031,
	"transactions_count":1
}`

var rawTransaction string = `{
	"blockHash":"000000000028e8b92df7010000f0eab92df7010000bd07f628fa7f000050e8b9",
	"blockIndex":0,
	"extra":{
		"nonce":[],
		"publicKey":"8bf1aecc80bf132fded9e1b4464042c04f4df223274e6d501cdbd419d55df9ae",
		"raw":"018bf1aecc80bf132fded9e1b4464042c04f4df223274e6d501cdbd419d55df9ae",
		"size":33
	},
	"fee":100000000000,
	"hash":"491857c7eb6276e0d872c1926d0c9863b39f207ae09d67b70b8bd648e16e82ab",
	"inBlockchain":false,
	"inputs":[
		{
				"data":{
					"input":{
							"amount":200000000,
							"k_image":"982f9b4063268ed8bc2e637db4e5e1195896b696a577f2278b5c2391b98efaee",
							"key_offsets":[
								25980,
								10137,
								48877
							]
					},
					"mixin":3,
					"outputs":[
							{
								"number":5,
								"transactionHash":"227fc58f1c1bdf89698c93a6280ce803cfe76606d6cfce7ef1eed5baf443b45c"
							},
							{
								"number":1,
								"transactionHash":"dbba432ed02c04540687ea05481ac782686e200efcc489241d08148d1c3fbd0b"
							},
							{
								"number":5,
								"transactionHash":"76d1c6a06c3049df25430ad2a50d42ec8947ac10b575487ba37f13ef26c17a12"
							}
					]
				},
				"type":"02"
		},
		{
				"data":{
					"input":{
							"amount":3000000000000,
							"k_image":"17275942b763fc3166ca9b86e6dfcab32ed5a556cc5b995a9c2bbf8a66d7e8d7",
							"key_offsets":[
								60338,
								22276,
								125912
							]
					},
					"mixin":3,
					"outputs":[
							{
								"number":10,
								"transactionHash":"b5c1b064dd0f8ee3ad60e5bb063276d8b49ea50c08487f290e5d525c9c77a4fb"
							},
							{
								"number":92,
								"transactionHash":"79993ad7abb2a03396129527571f2cd2cab468261e905992f6109d4dfbaebca7"
							},
							{
								"number":4,
								"transactionHash":"7ffdaefdc2f8d172633d4d22bfab044dcc2bfc44e9c60a4cd250e26af992bb0e"
							}
					]
				},
				"type":"02"
		}
	],
	"mixin":3,
	"outputs":[
		{
				"globalIndex":0,
				"output":{
					"amount":9,
					"target":{
							"data":{
								"key":"0d0de439543ff886bbf5d16e1cf4beb9e2163b1a50fadeb461aacf401ccf6351"
							},
							"type":"02"
					}
				}
		},
		{
				"globalIndex":0,
				"output":{
					"amount":1000000000000,
					"target":{
							"data":{
								"key":"df52982166519417263da2757d91d940bced321807d2e1e97df076cbfd790abf"
							},
							"type":"02"
					}
				}
		}
	],
	"paymentId":"0000000000000000000000000000000000000000000000000000000000000000",
	"signatures":[
		{
				"first":0,
				"second":"fc770928b3d5ac3a41b68c18e73b64a78e1a46259d1f4139a6a32b408f53720db01a9e4df8443f999bee42263ce5e88f31ab5fababbd69d0db9576775baf970e"
		},
		{
				"first":0,
				"second":"54006f8cbf961549c5cd8e4bb25293493efa6b20d843bd40afc921a0fcd1e40af2a58109b73f8dab78641894bd7a94266ed48b4a8d6355b2c766882617b97d0a"
		},
		{
				"first":24,
				"second":"3e4d53ae5b46592435d237981e6ad4cb5307e5e9e0d7be546e31bd99d6780a02c18c2847ecb5e8fb3920ae2264c113942b19a3b630157b8dedd8032117c9680c"
		}
	],
	"signaturesSize":25,
	"size":6577,
	"timestamp":1589473389,
	"totalInputsAmount":5609851418219,
	"totalOutputsAmount":5509851418219,
	"unlockTime":0,
	"version":1
}`

var rawShortTransaction string = `{
	"hash": "491857c7eb6276e0d872c1926d0c9863b39f207ae09d67b70b8bd648e16e82ab",
	"fee": 100000000000,
	"amount_out": 1000000000000,
	"size": 6577
}`

var rawMempoolTransaction string = `{
	"hash": "491857c7eb6276e0d872c1926d0c9863b39f207ae09d67b70b8bd648e16e82ab",
	"fee": 100000000000,
	"amount_out": 1000000000000,
	"size": 6577,
	"receive_time": 1589473389
}`

var rawStats string = `{
	"height": 123456,
  "already_generated_coins": 7822291504630269000,
  "transactions_count": 2,
  "block_size": 23716,
  "difficulty": 14377384679,
  "reward": 8507330332565,
  "timestamp": 1568118204
}`

type requestAssertion struct {
	method string
	params interface{}
	result string
}

var reqAssertion requestAssertion

func buildTestClient(t *testing.T) RPCKarbowanecd {
	assert := assert.New(t)

	testHTTPServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var request jsonrpc.RPCRequest

			data, _ := ioutil.ReadAll(r.Body)

			defer r.Body.Close()

			assert.Equal("application/json", r.Header.Get("Content-Type"), "Missing json header.")

			if err := json.Unmarshal(data, &request); err != nil {
				panic(err)
			}

			assert.Equal(reqAssertion.method, request.Method, "Method is not much")

			if reqAssertion.params != nil {
				assert.Equal(reqAssertion.params, request.Params, "Params are different from expected\n")
			} else {
				assert.Nil(request.Params)
			}

			fmt.Fprintf(w, "{\"jsonrpc\":\"2.0\",\"result\":%s}", reqAssertion.result)
		}),
	)

	// defer testHttpServer.Close()

	return NewClient(testHTTPServer.URL)
}

func toJSONResult(v interface{}) string {
	result, _ := json.Marshal(v)

	return string(result)
}

func assertBlock(t *testing.T, block Block) {
	assert.Equal(t, "4cc94394440e64461766d0e7c675498da934d70a63b143606246f0e6e26d1bde", block.Hash)
	assert.Equal(t, 7822291504630269000, block.AlreadyGeneratedCoins)
	assert.Equal(t, 842837, block.AlreadyGeneratedTransactions)
	assert.Equal(t, 8307330332565, block.BaseReward)
	assert.Equal(t, 23716, block.BlockSize)
	assert.Equal(t, 2277860049927920, block.CumulativeDifficulty)
	assert.Equal(t, 6, block.Depth)
	assert.Equal(t, 14377384679, block.Difficulty)
	assert.Equal(t, 1000000, block.EffectiveSizeMedian)
	assert.Equal(t, 399666, block.Index)
	assert.Equal(t, false, block.IsOrphaned)
	assert.Equal(t, 4, block.MajorVersion)
	assert.Equal(t, 0, block.MinorVersion)
	assert.Equal(t, 3653353692, block.Nonce)
	assert.Equal(t, float32(0), block.Penalty)
	assert.Equal(t, "72199371e92ee3dd95f0993f49a82d3dc7b7175dc004ad9c4437167a7138bd13", block.PrevBlockHash)
	assert.Equal(t, 8507330332565, block.Reward)
	assert.Equal(t, 532, block.SizeMedian)
	assert.Equal(t, 1568118204, block.Timestamp)
	assert.Equal(t, 200000000000, block.TotalFeeAmount)
	assert.Equal(t, 23640, block.TransactionsCumulativeSize)

	assert.Len(t, block.Transactions, 2)
}

func assertBlockHeader(t *testing.T, header BlockHeader) {
	assert.Equal(t, 1, header.Depth)
	assert.Equal(t, 65198, header.Difficulty)
	assert.Equal(t, "9a8be83...", header.Hash)
	assert.Equal(t, 123456, header.Height)
	assert.Equal(t, 1, header.MajorVersion)
	assert.Equal(t, 0, header.MinorVersion)
	assert.Equal(t, 2358499061, header.Nonce)
	assert.Equal(t, false, header.OrphanStatus)
	assert.Equal(t, "dde56b7e...", header.PrevHash)
	assert.Equal(t, 44090506423186, header.Reward)
	assert.Equal(t, 1356589561, header.Timestamp)
}

func assertTransaction(t *testing.T, trans Transaction) {
	assert.Equal(t, "000000000028e8b92df7010000f0eab92df7010000bd07f628fa7f000050e8b9", trans.BlockHash)
	assert.Equal(t, 0, trans.BlockIndex)
	assert.Equal(t, 100000000000, trans.Fee)
	assert.Equal(t, "491857c7eb6276e0d872c1926d0c9863b39f207ae09d67b70b8bd648e16e82ab", trans.Hash)
	assert.Equal(t, false, trans.InBlockchain)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", trans.PaymentID)
	assert.Equal(t, 25, trans.SignaturesSize)
	assert.Equal(t, 6577, trans.Size)
	assert.Equal(t, 3, trans.Mixin)
	assert.Equal(t, 1589473389, trans.Timestamp)
	assert.Equal(t, 5609851418219, trans.TotalInputsAmount)
	assert.Equal(t, 5509851418219, trans.TotalOutputsAmount)
	assert.Equal(t, 0, trans.UnlockTime)
	assert.Equal(t, 1, trans.Version)

	assert.Equal(t, []int{}, trans.Extra.Nonce)
	assert.Equal(t, "8bf1aecc80bf132fded9e1b4464042c04f4df223274e6d501cdbd419d55df9ae", trans.Extra.PublicKey)
	assert.Equal(t, "018bf1aecc80bf132fded9e1b4464042c04f4df223274e6d501cdbd419d55df9ae", trans.Extra.Raw)
	assert.Equal(t, 33, trans.Extra.Size)

	assert.Equal(t, 200000000, trans.Inputs[0].Data.Input.Amount)
	assert.Equal(t, "982f9b4063268ed8bc2e637db4e5e1195896b696a577f2278b5c2391b98efaee", trans.Inputs[0].Data.Input.KImage)
	assert.Equal(t, []int{25980, 10137, 48877}, trans.Inputs[0].Data.Input.KeyOffsets)
	assert.Equal(t, 3, trans.Inputs[0].Data.Mixin)
	assert.Equal(t, 5, trans.Inputs[0].Data.Outputs[0].Number)
	assert.Equal(t, "227fc58f1c1bdf89698c93a6280ce803cfe76606d6cfce7ef1eed5baf443b45c", trans.Inputs[0].Data.Outputs[0].TransactionHash)
	assert.Equal(t, 1, trans.Inputs[0].Data.Outputs[1].Number)
	assert.Equal(t, "dbba432ed02c04540687ea05481ac782686e200efcc489241d08148d1c3fbd0b", trans.Inputs[0].Data.Outputs[1].TransactionHash)
	assert.Equal(t, 5, trans.Inputs[0].Data.Outputs[2].Number)
	assert.Equal(t, "76d1c6a06c3049df25430ad2a50d42ec8947ac10b575487ba37f13ef26c17a12", trans.Inputs[0].Data.Outputs[2].TransactionHash)
	assert.Equal(t, "02", trans.Inputs[0].Type)

	assert.Equal(t, 3000000000000, trans.Inputs[1].Data.Input.Amount)
	assert.Equal(t, "17275942b763fc3166ca9b86e6dfcab32ed5a556cc5b995a9c2bbf8a66d7e8d7", trans.Inputs[1].Data.Input.KImage)
	assert.Equal(t, []int{60338, 22276, 125912}, trans.Inputs[1].Data.Input.KeyOffsets)
	assert.Equal(t, 3, trans.Inputs[1].Data.Mixin)
	assert.Equal(t, 10, trans.Inputs[1].Data.Outputs[0].Number)
	assert.Equal(t, "b5c1b064dd0f8ee3ad60e5bb063276d8b49ea50c08487f290e5d525c9c77a4fb", trans.Inputs[1].Data.Outputs[0].TransactionHash)
	assert.Equal(t, 92, trans.Inputs[1].Data.Outputs[1].Number)
	assert.Equal(t, "79993ad7abb2a03396129527571f2cd2cab468261e905992f6109d4dfbaebca7", trans.Inputs[1].Data.Outputs[1].TransactionHash)
	assert.Equal(t, 4, trans.Inputs[1].Data.Outputs[2].Number)
	assert.Equal(t, "7ffdaefdc2f8d172633d4d22bfab044dcc2bfc44e9c60a4cd250e26af992bb0e", trans.Inputs[1].Data.Outputs[2].TransactionHash)
	assert.Equal(t, "02", trans.Inputs[1].Type)

	assert.Equal(t, 0, trans.Outputs[0].GlobalIndex)
	assert.Equal(t, 9, trans.Outputs[0].Output.Amount)
	assert.Equal(t, "0d0de439543ff886bbf5d16e1cf4beb9e2163b1a50fadeb461aacf401ccf6351", trans.Outputs[0].Output.Target.Data.Key)
	assert.Equal(t, "02", trans.Outputs[0].Output.Target.Type)

	assert.Equal(t, 0, trans.Outputs[1].GlobalIndex)
	assert.Equal(t, 1000000000000, trans.Outputs[1].Output.Amount)
	assert.Equal(t, "df52982166519417263da2757d91d940bced321807d2e1e97df076cbfd790abf", trans.Outputs[1].Output.Target.Data.Key)
	assert.Equal(t, "02", trans.Outputs[1].Output.Target.Type)

	assert.Equal(t, 0, trans.Signatures[0].First)
	assert.Equal(t, "fc770928b3d5ac3a41b68c18e73b64a78e1a46259d1f4139a6a32b408f53720db01a9e4df8443f999bee42263ce5e88f31ab5fababbd69d0db9576775baf970e", trans.Signatures[0].Second)
	assert.Equal(t, 0, trans.Signatures[1].First)
	assert.Equal(t, "54006f8cbf961549c5cd8e4bb25293493efa6b20d843bd40afc921a0fcd1e40af2a58109b73f8dab78641894bd7a94266ed48b4a8d6355b2c766882617b97d0a", trans.Signatures[1].Second)
	assert.Equal(t, 24, trans.Signatures[2].First)
	assert.Equal(t, "3e4d53ae5b46592435d237981e6ad4cb5307e5e9e0d7be546e31bd99d6780a02c18c2847ecb5e8fb3920ae2264c113942b19a3b630157b8dedd8032117c9680c", trans.Signatures[2].Second)
}

func assertMempoolTransaction(t *testing.T, trans MempoolTransaction) {
	assert.Equal(t, "491857c7eb6276e0d872c1926d0c9863b39f207ae09d67b70b8bd648e16e82ab", trans.Hash)
	assert.Equal(t, 100000000000, trans.Fee)
	assert.Equal(t, 1000000000000, trans.AmountOut)
	assert.Equal(t, 6577, trans.Size)
	assert.Equal(t, 1589473389, trans.ReceiveTime)
}

func assertShortTransaction(t *testing.T, trans ShortTransaction) {
	assert.Equal(t, "491857c7eb6276e0d872c1926d0c9863b39f207ae09d67b70b8bd648e16e82ab", trans.Hash)
	assert.Equal(t, 100000000000, trans.Fee)
	assert.Equal(t, 1000000000000, trans.AmountOut)
	assert.Equal(t, 6577, trans.Size)
}

func assertStats(t *testing.T, stats BlockStats) {
	assert.Equal(t, 123456, stats.Height)
	assert.Equal(t, 7822291504630269000, stats.AlreadyGeneratedCoins)
	assert.Equal(t, 2, stats.TransactionsCount)
	assert.Equal(t, 23716, stats.BlockSize)
	assert.Equal(t, 14377384679, stats.Difficulty)
	assert.Equal(t, 8507330332565, stats.Reward)
	assert.Equal(t, 1568118204, stats.Timestamp)
}

func TestGetBlockCount(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "getblockcount",
		params: map[string]interface{}{},
		result: toJSONResult(map[string]interface{}{
			"count":  123456,
			"status": "OK",
		}),
	}

	count, err := testClient.GetBlockCount()

	assert.Nil(t, err)
	assert.Equal(t, 123456, count)
}

func TestGetBlockHash(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "getblockhash",
		params: []interface{}([]interface{}{float64(123456)}),
		// params: map[string]interface{}{
		// 	"height": float64(123456),
		// },
		result: "\"test_hash\"",
	}

	hash, err := testClient.GetBlockHash(123456)

	assert.Nil(t, err)
	assert.Equal(t, "test_hash", hash)
}

func TestGetBlockTemplate(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "getblocktemplate",
		params: map[string]interface{}{
			"reserve_size":   float64(123),
			"wallet_address": testAddress,
		},
		result: toJSONResult(map[string]interface{}{
			"blocktemplate_blob": "0100de...",
			"blockhashing_blob":  "0400fba1a6f805...",
			"difficulty":         65563,
			"height":             123456,
			"reserved_offset":    395,
			"status":             "OK",
		}),
	}

	blockTemplate, err := testClient.GetBlockTemplate(123, testAddress)

	assert.Nil(t, err)
	assert.Equal(t, "0100de...", blockTemplate.BlockTemplateBlob)
	assert.Equal(t, "0400fba1a6f805...", blockTemplate.BlockHashingBlob)
	assert.Equal(t, 65563, blockTemplate.Difficulty)
	assert.Equal(t, 123456, blockTemplate.Height)
	assert.Equal(t, 395, blockTemplate.ReservedOffset)
}

func TestGetBlockHeaderByHash(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "getblockheaderbyhash",
		params: map[string]interface{}{
			"hash": "0100de...",
		},
		result: toJSONResult(map[string]interface{}{
			"block_header": testBlockHeader,
			"status":       "OK",
		}),
	}

	header, err := testClient.GetBlockHeaderByHash("0100de...")

	assert.Nil(t, err)

	assert.Equal(t, 1, header.Depth)
	assert.Equal(t, 65198, header.Difficulty)
	assert.Equal(t, "9a8be83...", header.Hash)
	assert.Equal(t, 123456, header.Height)
	assert.Equal(t, 1, header.MajorVersion)
	assert.Equal(t, 0, header.MinorVersion)
	assert.Equal(t, 2358499061, header.Nonce)
	assert.Equal(t, false, header.OrphanStatus)
	assert.Equal(t, "dde56b7e...", header.PrevHash)
	assert.Equal(t, 44090506423186, header.Reward)
	assert.Equal(t, 1356589561, header.Timestamp)
}

func TestGetBlockHeaderByHeight(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "getblockheaderbyheight",
		params: map[string]interface{}{
			"height": float64(123456),
		},
		result: toJSONResult(map[string]interface{}{
			"block_header": testBlockHeader,
			"status":       "OK",
		}),
	}

	header, err := testClient.GetBlockHeaderByHeight(123456)

	assert.Nil(t, err)
	assertBlockHeader(t, header)
}

func TestGetBlockTimestamp(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "getblocktimestamp",
		params: map[string]interface{}{
			"height": float64(123456),
		},
		result: toJSONResult(map[string]interface{}{
			"timestamp": 1356589561,
			"status":    "OK",
		}),
	}

	timestamp, err := testClient.GetBlockTimestamp(123456)

	assert.Nil(t, err)
	assert.Equal(t, 1356589561, timestamp)
}

func TestGetBlockByHash(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "getblockbyhash",
		params: map[string]interface{}{
			"hash": "test_hash",
		},
		result: `{
			"block":` + rawBlock + `,
			"status": "OK"
		}`,
	}

	block, err := testClient.GetBlockByHash("test_hash")

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)

	assertBlock(t, block)
}

func TestGetBlockByHeight(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "getblockbyheight",
		params: map[string]interface{}{
			"blockHeight": float64(123456),
		},
		result: `{
			"block":` + rawBlock + `,
			"status": "OK"
		}`,
	}

	block, err := testClient.GetBlockByHeight(123456)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)

	assertBlock(t, block)
}

func TestGetBlocksByHeighs(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "getblocksbyheights",
		params: map[string]interface{}{
			"blockHeights": []interface{}{float64(123456), float64(123457), float64(123458)},
		},
		result: `{
			"blocks":[` + rawBlock + `,` + rawBlock + `,` + rawBlock + `],
			"status": "OK"
		}`,
	}

	blocks, err := testClient.GetBlocksByHeights([]int{123456, 123457, 123458})

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.Len(t, blocks, 3)
	assertBlock(t, blocks[0])
	assertBlock(t, blocks[1])
	assertBlock(t, blocks[2])
}

func TestGetBlocksByHashes(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "getblocksbyhashes",
		params: map[string]interface{}{
			"blockHashes": []interface{}{
				"4cc94394440e64461766d0e7c675498da934d70a63b143606246f0e6e26d1bde",
				"4cc94394440e64461766d0e7c675498da934d70a63b143606246f0e6e26d1bde",
				"4cc94394440e64461766d0e7c675498da934d70a63b143606246f0e6e26d1bde",
			},
		},
		result: `{
			"blocks":[` + rawBlock + `,` + rawBlock + `,` + rawBlock + `],
			"status": "OK"
		}`,
	}

	blocks, err := testClient.GetBlocksByHashes([]string{
		"4cc94394440e64461766d0e7c675498da934d70a63b143606246f0e6e26d1bde",
		"4cc94394440e64461766d0e7c675498da934d70a63b143606246f0e6e26d1bde",
		"4cc94394440e64461766d0e7c675498da934d70a63b143606246f0e6e26d1bde",
	})

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.Len(t, blocks, 3)
	assertBlock(t, blocks[0])
	assertBlock(t, blocks[1])
	assertBlock(t, blocks[2])
}

func TestGetBlocksHashesByTimestamps(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "getblockshashesbytimestamps",
		params: map[string]interface{}{
			"timestampBegin": float64(123456),
			"timestampEnd":   float64(133333),
			"limit":          float64(3),
		},
		result: `{
			"blockHashes":["th1", "th2", "th3"],
			"count": 5,
			"status": "OK"
		}`,
	}

	hashes, count, err := testClient.GetBlocksHashesByTimestamps(123456, 133333, 3)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.Len(t, hashes, 3)
	assert.Equal(t, 5, count)
	assert.Equal(t, "th1", hashes[0])
	assert.Equal(t, "th2", hashes[1])
	assert.Equal(t, "th3", hashes[2])
}

func TestGetBlocksList(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "getblockslist",
		params: map[string]interface{}{
			"height": float64(123456),
			"count":  float64(3),
		},
		result: `{
			"blocks":[` + rawShortBlock + `,` + rawShortBlock + `,` + rawShortBlock + `]
		}`,
	}

	blocks, err := testClient.GetBlocksList(123456, 3)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.Len(t, blocks, 3)
}

func TestGetAltBlocksList(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "getaltblockslist",
		// params: map[string]interface{}{},
		params: nil,
		result: `{
			"alt_blocks":[` + rawShortBlock + `,` + rawShortBlock + `,` + rawShortBlock + `]
		}`,
	}

	blocks, err := testClient.GetAltBlocksList()

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.Len(t, blocks, 3)
}

func TestGetLastBlockHeader(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "getlastblockheader",
		params: nil,
		result: toJSONResult(map[string]interface{}{
			"block_header": testBlockHeader,
			"status":       "OK",
		}),
	}

	header, err := testClient.GetLastBlockHeader()

	assert.Nil(t, err)
	assertBlockHeader(t, header)
}

func TestGetTransaction(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "gettransaction",
		params: map[string]interface{}{
			"hash": "491857c7eb6276e0d872c1926d0c9863b39f207ae09d67b70b8bd648e16e82ab",
		},
		result: `{
			"transaction":` + rawTransaction + `,
			"status": "OK"
		}`,
	}

	transaction, err := testClient.GetTransaction("491857c7eb6276e0d872c1926d0c9863b39f207ae09d67b70b8bd648e16e82ab")

	assert.Nil(t, err)
	assertTransaction(t, transaction)
}

func TestGetTransactionsPool(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "gettransactionspool",
		params: nil,
		result: `{
			"transactions": [` + rawMempoolTransaction + `,` + rawMempoolTransaction + `],
			"status":       "OK"
		}`,
	}

	mempoolTranscations, err := testClient.GetTransactionsPool()

	assert.Nil(t, err)

	assert.Len(t, mempoolTranscations, 2)
	assertMempoolTransaction(t, mempoolTranscations[0])
	assertMempoolTransaction(t, mempoolTranscations[1])
}

func TestGetTransactionByPaymentID(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "gettransactionsbypaymentid",
		params: map[string]interface{}{
			"payment_id": "pid",
		},
		result: `{
			"transactions": [` + rawShortTransaction + `,` + rawShortTransaction + `],
			"status":       "OK"
		}`,
	}

	transactions, err := testClient.GetTransactionsByPaymentID("pid")

	assert.Nil(t, err)

	assert.Len(t, transactions, 2)
	assertShortTransaction(t, transactions[0])
	assertShortTransaction(t, transactions[1])
}

func TestGetTransactionHashesByPaymentID(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "gettransactionhashesbypaymentid",
		params: map[string]interface{}{
			"paymentId": "pid",
		},
		result: `{
			"transactionHashes": ["thash1","thash2"],
			"status":       "OK"
		}`,
	}

	hashes, err := testClient.GetTransactionHashesByPaymentID("pid")

	assert.Nil(t, err)
	assert.Len(t, hashes, 2)
	assert.Equal(t, "thash1", hashes[0])
	assert.Equal(t, "thash2", hashes[1])
}

func TestGetTransactionsByHashes(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "gettransactionsbyhashes",
		params: map[string]interface{}{
			"transactionHashes": []interface{}{"thash1", "thash2"},
		},
		result: `{
			"transactions": [` + rawTransaction + `,` + rawTransaction + `],
			"status":       "OK"
		}`,
	}

	transactions, err := testClient.GetTransactionsByHashes([]string{"thash1", "thash2"})

	assert.Nil(t, err)
	assert.Len(t, transactions, 2)
	assertTransaction(t, transactions[0])
	assertTransaction(t, transactions[1])
}

func TestGetCurrencyID(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "getcurrencyid",
		params: nil,
		result: `{
			"currency_id_blob": "curr_blob"
		}`,
	}

	id, err := testClient.GetCurrencyID()

	assert.Nil(t, err)
	assert.Equal(t, "curr_blob", id)
}

func TestGetStatsByHeights(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "getstatsbyheights",
		params: map[string]interface{}{
			"heights": []interface{}{float64(123456), float64(123457)},
		},
		result: `{
			"stats": [` + rawStats + `,` + rawStats + `],
			"status":       "OK"
		}`,
	}

	stats, err := testClient.GetStatsByHeights([]int{123456, 123457})

	assert.Nil(t, err)
	assert.Len(t, stats, 2)
	assertStats(t, stats[0])
	assertStats(t, stats[1])
}

func TestGetStatsByRange(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "getstatsinrange",
		params: map[string]interface{}{
			"start_height": float64(123456),
			"end_height":   float64(123458),
		},
		result: `{
			"stats": [` + rawStats + `,` + rawStats + `,` + rawStats + `],
			"status":       "OK"
		}`,
	}

	stats, err := testClient.GetStatsInRange(123456, 123458)

	assert.Nil(t, err)
	assert.Len(t, stats, 3)
	assertStats(t, stats[0])
	assertStats(t, stats[1])
	assertStats(t, stats[2])
}

func TestCheckTransactionKey(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "checktransactionkey",
		params: map[string]interface{}{
			"transaction_id":  "494862266B7E511A4F901E7A0EFFB27B012ABA4447D4B0666A010781B7A6FBE5",
			"address":         "KhnaNiHdEdWcLh8VseUa8Z2fv2e6w2rvWgYySWZaEUMSCHFqhJyhWXY6tUeYsqSZ31QANbaZP6Gm7byZomD4xsc9QMye6ZX",
			"transaction_key": "43210C433EAE73271779652F67DFF20B03591776FE76FEA4FD55D69F237CBC0D",
		},
		result: `{
      "amount":300000000000000,
      "outputs":[
         {
            "amount":300000000000000,
            "target":{
               "data":{
                  "key":"a5bfae862500ab22e600ff90fb7052960fa2e0abb4d129c6c4b494cc4be1337f"
               },
               "type":"02"
            }
         }
      ],
      "status":"OK"
   }`,
	}

	check, err := testClient.CheckTransactionKey(
		"494862266B7E511A4F901E7A0EFFB27B012ABA4447D4B0666A010781B7A6FBE5",
		"43210C433EAE73271779652F67DFF20B03591776FE76FEA4FD55D69F237CBC0D",
		"KhnaNiHdEdWcLh8VseUa8Z2fv2e6w2rvWgYySWZaEUMSCHFqhJyhWXY6tUeYsqSZ31QANbaZP6Gm7byZomD4xsc9QMye6ZX",
	)

	assert.Nil(t, err)

	assert.Equal(t, 300000000000000, check.Amount)
	assert.Equal(t, 300000000000000, check.Outputs[0].Amount)
	assert.Equal(t, "a5bfae862500ab22e600ff90fb7052960fa2e0abb4d129c6c4b494cc4be1337f", check.Outputs[0].Target.Data.Key)
	assert.Equal(t, "02", check.Outputs[0].Target.Type)

}

func TestCheckTransactionByViewKey(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "checktransactionbyviewkey",
		params: map[string]interface{}{
			"transaction_id": "462266B7E511A4F901E7A0EFFB27B012ABA4948447D4B0666A010781B7A6FBE5",
			"view_key":       "4bcd75aef0eaf0c1e34caf024cc1f499a7656083e4921e58552c3c1254e99209",
			"address":        "Kiwr4Aajn7QefxbAEvHSAcTrVbhzYukfmbDvwWkrDhjFXe1FVUa5ggxCrEv4w2zvaiVZgoF4v7b3cNAbaU3LKQGS9EU9KAd",
		},
		result: `{
      "amount":599900000600000,
      "confirmations":5018,
      "outputs":[
         {
            "amount":600000,
            "target":{
               "data":{
                  "key":"09f79dc4a9c0cc6e32be6cb32a70cc4a4f9f09185627d234470d830004e665eb"
               },
               "type":"02"
            }
         },
         {
            "amount":500000000000000,
            "target":{
               "data":{
                  "key":"21d49839518f2d07524e6710093cdbcedfe4076265f851d117925d18e01ee628"
               },
               "type":"02"
            }
         }
			],
			"status":"OK"
    }`,
	}

	check, err := testClient.CheckTransactionByViewKey(
		"462266B7E511A4F901E7A0EFFB27B012ABA4948447D4B0666A010781B7A6FBE5",
		"4bcd75aef0eaf0c1e34caf024cc1f499a7656083e4921e58552c3c1254e99209",
		"Kiwr4Aajn7QefxbAEvHSAcTrVbhzYukfmbDvwWkrDhjFXe1FVUa5ggxCrEv4w2zvaiVZgoF4v7b3cNAbaU3LKQGS9EU9KAd",
	)

	assert.Nil(t, err)

	assert.Equal(t, 599900000600000, check.Amount)
	assert.Equal(t, 5018, check.Confirmations)

	assert.Equal(t, 600000, check.Outputs[0].Amount)
	assert.Equal(t, "09f79dc4a9c0cc6e32be6cb32a70cc4a4f9f09185627d234470d830004e665eb", check.Outputs[0].Target.Data.Key)
	assert.Equal(t, "02", check.Outputs[0].Target.Type)

	assert.Equal(t, 500000000000000, check.Outputs[1].Amount)
	assert.Equal(t, "21d49839518f2d07524e6710093cdbcedfe4076265f851d117925d18e01ee628", check.Outputs[1].Target.Data.Key)
	assert.Equal(t, "02", check.Outputs[1].Target.Type)

}

func TestCheckTransactionProof(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "checktransactionproof",
		params: map[string]interface{}{
			"transaction_id":      "2196e40ff4bcd7711f57d7ccd2aab67f8dddf6015ad72ae2b4d63621521a94a9",
			"destination_address": "Kiwr4Aajn7QefxbAEvHSAcTrVbhzYukfmbDvwWkrDhjFXe1FVUa5ggxCrEv4w2zvaiVZgoF4v7b3cNAbaU3LKQGS9EU9KAd",
			"signature":           "ProofR3Y7go5QVgEx8zvs4LWQ5ghzRifmUtLDXYhEMm2Ku9idhg8S3NEM4eeWDZbGEFCdn6H4GoyZcnQ2FHWb8hteHD3BkGBYDbZ1hA7xoRGY9DX9xCSXSQ2Xb7HDXARjX4YjGfyAWgGZ9f",
		},
		result: `{
      "confirmations":23294,
      "outputs":[
         {
            "amount":100000000000,
            "target":{
               "data":{
                  "key":"de5a651865919d253d08b88bf78f3923f02f9f4002219f23f2bfdc6fac5a9ff6"
               },
               "type":"02"
            }
         },
         {
            "amount":2000000000000,
            "target":{
               "data":{
                  "key":"5f4feea8da71e322dd60c1057ffe80d288322f211ab30e817889d71bc134d040"
               },
               "type":"02"
            }
         }
      ],
      "received_amount":2100000000000,
      "signature_valid":true,
      "status":"OK"
   }`,
	}

	check, err := testClient.CheckTransactionProof(
		"2196e40ff4bcd7711f57d7ccd2aab67f8dddf6015ad72ae2b4d63621521a94a9",
		"ProofR3Y7go5QVgEx8zvs4LWQ5ghzRifmUtLDXYhEMm2Ku9idhg8S3NEM4eeWDZbGEFCdn6H4GoyZcnQ2FHWb8hteHD3BkGBYDbZ1hA7xoRGY9DX9xCSXSQ2Xb7HDXARjX4YjGfyAWgGZ9f",
		"Kiwr4Aajn7QefxbAEvHSAcTrVbhzYukfmbDvwWkrDhjFXe1FVUa5ggxCrEv4w2zvaiVZgoF4v7b3cNAbaU3LKQGS9EU9KAd",
	)

	assert.Nil(t, err)

	assert.Equal(t, 23294, check.Confirmations)
	assert.Equal(t, 2100000000000, check.ReceivedAmount)
	assert.Equal(t, true, check.SignatureValid)

	assert.Equal(t, 100000000000, check.Outputs[0].Amount)
	assert.Equal(t, "de5a651865919d253d08b88bf78f3923f02f9f4002219f23f2bfdc6fac5a9ff6", check.Outputs[0].Target.Data.Key)
	assert.Equal(t, "02", check.Outputs[0].Target.Type)

	assert.Equal(t, 2000000000000, check.Outputs[1].Amount)
	assert.Equal(t, "5f4feea8da71e322dd60c1057ffe80d288322f211ab30e817889d71bc134d040", check.Outputs[1].Target.Data.Key)
	assert.Equal(t, "02", check.Outputs[1].Target.Type)
}

func TestCheckReserveProof(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "checkreserveproof",
		params: map[string]interface{}{
			"height":    float64(489000),
			"address":   "Kiwr4Aajn7QefxbAEvHSAcTrVbhzYukfmbDvwWkrDhjFXe1FVUa5ggxCrEv4w2zvaiVZgoF4v7b3cNAbaU3LKQGS9EU9KAd",
			"message":   "this is challenge message",
			"signature": "RsrvPrf......",
		},
		result: `{
      "good":true,
      "locked":0,
      "spent":0,
      "total":103143982017908
   }`,
	}

	check, err := testClient.CheckReserveProof(
		"Kiwr4Aajn7QefxbAEvHSAcTrVbhzYukfmbDvwWkrDhjFXe1FVUa5ggxCrEv4w2zvaiVZgoF4v7b3cNAbaU3LKQGS9EU9KAd",
		"RsrvPrf......",
		"this is challenge message",
		489000,
	)

	assert.Nil(t, err)

	assert.Equal(t, true, check.Good)
	assert.Equal(t, 0, check.Locked)
	assert.Equal(t, 0, check.Spent)
	assert.Equal(t, 103143982017908, check.Total)
}

func TestValidateAddress(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "validateaddress",
		params: map[string]interface{}{
			"address": "Kd79ZmERRBx3t5uuERi6ggTMMeTBDTLY21cvkEFtE29GY8FJuwQHmWyRxYbPxYBu8S8a7wzxhg3tJfbE27hYGYtbD4mGiSs",
		},
		result: `{
      "address":"Kd79ZmERRBx3t5uuERi6ggTMMeTBDTLY21cvkEFtE29GY8FJuwQHmWyRxYbPxYBu8S8a7wzxhg3tJfbE27hYGYtbD4mGiSs",
      "is_valid":true,
      "spend_public_key":"56359b8d0f44eb113916a3dfdcceb99d8ac8b904d7e8c303b40b6e8356c3bbba",
      "status":"OK",
      "view_public_key":"1578f4ff1adf2a95364f118ef2f05f2d43a50693d3491fe6b709ff7db9ea546a"
   }`,
	}

	validation, err := testClient.ValidateAddress(
		"Kd79ZmERRBx3t5uuERi6ggTMMeTBDTLY21cvkEFtE29GY8FJuwQHmWyRxYbPxYBu8S8a7wzxhg3tJfbE27hYGYtbD4mGiSs",
	)

	assert.Nil(t, err)

	assert.Equal(t, true, validation.IsValid)
	assert.Equal(t, "Kd79ZmERRBx3t5uuERi6ggTMMeTBDTLY21cvkEFtE29GY8FJuwQHmWyRxYbPxYBu8S8a7wzxhg3tJfbE27hYGYtbD4mGiSs", validation.Address)
	assert.Equal(t, "56359b8d0f44eb113916a3dfdcceb99d8ac8b904d7e8c303b40b6e8356c3bbba", validation.SpendPublicKey)
	assert.Equal(t, "1578f4ff1adf2a95364f118ef2f05f2d43a50693d3491fe6b709ff7db9ea546a", validation.ViewPublicKey)
}

func TestVerifyMessage(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "verifymessage",
		params: map[string]interface{}{
			"message":   "test",
			"address":   "Kiwr4Aajn7QefxbAEvHSAcTrVbhzYukfmbDvwWkrDhjFXe1FVUa5ggxCrEv4w2zvaiVZgoF4v7b3cNAbaU3LKQGS9EU9KAd",
			"signature": "SigV1RwnRhiLR8MU2S55xj9kEt119U4JcEzbToamAbsZ9qzRD3tDxUEBCCWMNDKCeGqDh25g6Wguq9Fio4SteFfwkUSZP",
		},
		result: `{
      "sig_valid":true,
      "status":"OK"
   }`,
	}

	valid, err := testClient.VerifyMessage(
		"test",
		"Kiwr4Aajn7QefxbAEvHSAcTrVbhzYukfmbDvwWkrDhjFXe1FVUa5ggxCrEv4w2zvaiVZgoF4v7b3cNAbaU3LKQGS9EU9KAd",
		"SigV1RwnRhiLR8MU2S55xj9kEt119U4JcEzbToamAbsZ9qzRD3tDxUEBCCWMNDKCeGqDh25g6Wguq9Fio4SteFfwkUSZP",
	)

	assert.Nil(t, err)
	assert.True(t, valid)
}

func TestSubmitBlock(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "submitblock",
		params: map[string]interface{}{
			"block_blob": "test",
		},
		result: `{
      "status":"OK"
   }`,
	}

	err := testClient.SubmitBlock("test")

	assert.Nil(t, err)
}
