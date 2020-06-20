package walletd_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/r3volut1oner/go-karbo/client/walletd"
	"github.com/ybbus/jsonrpc"

	"github.com/stretchr/testify/assert"
)

type requestAssertion struct {
	method string
	params interface{}
	result string
}

var httpServer *httptest.Server
var reqAssertion requestAssertion

func toJSONResult(v interface{}) string {
	result, _ := json.Marshal(v)

	return string(result)
}

var testAddress = "KfXkT5VmdqmA7bWqSH37p87hSXBdTpTogN4mGHPARUSJaLse6jbXaVbVkLs3DwcmuD88xfu835Zvh6qBPCUXw6CHK8koDCt"

var testTransactionResponse = map[string]interface{}{
	"fee":             -70368475742208,
	"extra":           "0127cea59bfadc49aa02ed4a225936671e55607b5241621abca2a5e14405906dbb",
	"timestamp":       1446029698,
	"blockIndex":      1,
	"state":           0,
	"transactionHash": "06ec210a8359f253f8b2160a0d6040cf89f2a05a553aaa577b7f508ee5d831f9",
	"amount":          70368475742208,
	"unlockTime":      11,
	"transfers": []map[string]interface{}{
		{
			"amount":  70368475742208,
			"type":    0,
			"address": "KfXkT5VmdqmA7bWqSH37p87hSXBdTpTogN4mGHPARUSJaLse6jbXaVbVkLs3DwcmuD88xfu835Zvh6qBPCUXw6CHK8koDCt",
		},
	},
	"paymentId": "",
	"isBase":    true,
}

func buildTestClient(t *testing.T) walletd.RPCWalletd {
	assert := assert.New(t)

	testHttpServer := httptest.NewServer(
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

	return walletd.NewClient(testHttpServer.URL)
}

func TestSave(t *testing.T) {
	testClient := buildTestClient(t)
	reqAssertion = requestAssertion{
		method: "save",
		params: nil,
		result: "{}",
	}

	err := testClient.Save()

	if err != nil {
		t.Error(err)
	}

	assert.Nil(t, err)
}

func TestReset(t *testing.T) {
	reqAssertion = requestAssertion{
		method: "reset",
		result: "{}",
	}

	testClient := buildTestClient(t)
	err := testClient.Reset(nil)

	if err != nil {
		t.Error(err)
		return
	}

	reqAssertion = requestAssertion{
		method: "reset",
		params: map[string]interface{}{
			"viewSecretKey": "secret_key_test",
		},
		result: "{}",
	}

	err = testClient.Reset(&walletd.ReqReset{
		ViewSecretKey: "secret_key_test",
	})

	if err != nil {
		t.Error(err)
		return
	}

	reqAssertion = requestAssertion{
		method: "reset",
		params: map[string]interface{}{
			"scanHeight": float64(10000),
		},
		result: "{}",
	}

	err = testClient.Reset(&walletd.ReqReset{
		ScanHeight: 10000,
	})

	if err != nil {
		t.Error(err)
		return
	}
}

func TestExport(t *testing.T) {
	reqAssertion = requestAssertion{
		method: "export",
		params: map[string]interface{}{
			"fileName": "container_file",
		},
		result: "{}",
	}

	testClient := buildTestClient(t)
	err := testClient.Export(&walletd.ReqExport{
		FileName: "container_file",
	})

	if err != nil {
		t.Error(err)
		return
	}
}

func TestCreateAddress(t *testing.T) {
	reqAssertion = requestAssertion{
		method: "createAddress",
		result: toJSONResult(map[string]string{
			"address": testAddress,
		}),
	}

	testClient := buildTestClient(t)
	response, err := testClient.CreateAddress(nil)

	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, &walletd.ResCreateAddress{testAddress}, response)
}

func TestCreateAddressList(t *testing.T) {
	reqAssertion = requestAssertion{
		method: "createAddressList",
		params: map[string]interface{}{
			"spendSecretKeys": []interface{}{
				"testSpendKey1",
				"testSpendKey2",
			},
		},
		result: toJSONResult(map[string]interface{}{
			"addresses": []string{
				"testAddress1",
				"testAddress2",
			},
		}),
	}

	testClient := buildTestClient(t)
	response, err := testClient.CreateAddressList(&walletd.ReqCreateAddressList{
		SpendSecretKeys: []string{
			"testSpendKey1",
			"testSpendKey2",
		},
	})

	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(
		t,
		&walletd.ResCreateAddressList{
			Addresses: []string{
				"testAddress1",
				"testAddress2",
			},
		},
		response,
	)
}

func TestDeleteAddress(t *testing.T) {
	reqAssertion = requestAssertion{
		method: "deleteAddress",
		params: map[string]interface{}{
			"address": testAddress,
		},
		result: "{}",
	}

	testClient := buildTestClient(t)
	err := testClient.DeleteAddress(testAddress)

	assert.Nil(t, err)
}

func TestGetSpendKeys(t *testing.T) {
	reqAssertion = requestAssertion{
		method: "getSpendKeys",
		params: map[string]interface{}{
			"address": testAddress,
		},
		result: toJSONResult(map[string]interface{}{
			"spendSecretKey": "testSecretKey",
			"spendPublicKey": "testPublicKey",
		}),
	}

	testClient := buildTestClient(t)
	response, err := testClient.GetSpendKeys(testAddress)

	assert.Nil(t, err)
	assert.Equal(t, "testSecretKey", response.SpendSecretKey)
	assert.Equal(t, "testPublicKey", response.SpendPublicKey)
}

func TestGetBalance(t *testing.T) {
	reqAssertion = requestAssertion{
		method: "getBalance",
		params: map[string]interface{}{
			"address": testAddress,
		},
		result: toJSONResult(map[string]uint64{
			"availableBalance": 1000000000,
			"lockedAmount":     0,
		}),
	}

	testClient := buildTestClient(t)
	response, err := testClient.GetBalance(testAddress)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, &walletd.ResBalance{
		AvailableBalance: 1000000000,
		LockedAmount:     0,
	}, response)
}

func TestGetBalanceDefault(t *testing.T) {
	reqAssertion = requestAssertion{
		method: "getBalance",
		result: toJSONResult(map[string]uint64{
			"availableBalance": 100000000,
			"lockedAmount":     200000000,
		}),
	}

	testClient := buildTestClient(t)
	response, err := testClient.GetBalance("")

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, &walletd.ResBalance{
		AvailableBalance: 100000000,
		LockedAmount:     200000000,
	}, response)
}

func TestGetBlockHashes(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "getBlockHashes",
		params: map[string]interface{}{
			"blockCount":      float64(11),
			"firstBlockIndex": float64(0),
		},
		result: toJSONResult(map[string][]string{
			"blockHashes": {
				"8a6f1cb7ed7a9db4751d7b283a0482baff20567173dbfae136c9bceb188e51c4",
				"657cd1c33df7f4250d581c97db665cb4a1856ebfadd6efabaeab745c2c76b1be",
				"21047174f74576b6722e72646d7bd553e17d7c9f07fef05151bb1f9df7ed9dd8",
			},
		}),
	}

	response, err := tclient.GetBlockHashes(&walletd.ReqGetBlockHashes{
		BlockCount:      11,
		FirstBlockIndex: 0,
	})

	assert.Nil(t, err)
	assert.Equal(t, "8a6f1cb7ed7a9db4751d7b283a0482baff20567173dbfae136c9bceb188e51c4", response.BlockHashes[0])
	assert.Equal(t, "657cd1c33df7f4250d581c97db665cb4a1856ebfadd6efabaeab745c2c76b1be", response.BlockHashes[1])
	assert.Equal(t, "21047174f74576b6722e72646d7bd553e17d7c9f07fef05151bb1f9df7ed9dd8", response.BlockHashes[2])
}

func TestGetTransactionHashes(t *testing.T) {
	tclient := buildTestClient(t)
	blockHash := "test_block_hash"
	firstBlockIndex := 10

	_, err := tclient.GetTransactionHashes(&walletd.ReqGetTransactionHashes{
		BlockHash:       blockHash,
		FirstBlockIndex: firstBlockIndex,
	})

	assert.Equal(t, "one of 'BlockHash' or 'FirstBlockIndex' must be provided", err.Error())

	reqAssertion = requestAssertion{
		method: "getTransactionHashes",
		params: map[string]interface{}{
			"addresses":  []interface{}{testAddress, testAddress},
			"blockHash":  blockHash,
			"paymentId":  "test_payment_id",
			"blockCount": float64(0),
		},
		result: toJSONResult(map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"transactionHashes": []string{"trans1", "trans2"},
					"blockHash":         "testBlockHash",
				},
			},
		}),
	}

	response, err := tclient.GetTransactionHashes(&walletd.ReqGetTransactionHashes{
		Addresses: []string{testAddress, testAddress},
		PaymentID: "test_payment_id",
		BlockHash: blockHash,
	})

	assert.Equal(t, "testBlockHash", response.Items[0].BlockHash)
	assert.Equal(t, "trans1", response.Items[0].TransactionHashes[0])
	assert.Equal(t, "trans2", response.Items[0].TransactionHashes[1])
}

func TestGetTransactions(t *testing.T) {
	tclient := buildTestClient(t)
	blockHash := "test_block_hash"
	firstBlockIndex := 10

	assert := assert.New(t)

	_, err := tclient.GetTransactionHashes(&walletd.ReqGetTransactionHashes{
		BlockHash:       blockHash,
		FirstBlockIndex: firstBlockIndex,
	})

	assert.Equal("one of 'BlockHash' or 'FirstBlockIndex' must be provided", err.Error())

	reqAssertion = requestAssertion{
		method: "getTransactions",
		params: map[string]interface{}{
			"addresses":  []interface{}{testAddress, testAddress},
			"blockHash":  blockHash,
			"paymentId":  "test_payment_id",
			"blockCount": float64(0),
		},
		result: toJSONResult(map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"blockHash": "01bd06ca731914f27e143bbb902ce0bc05bff13d76faa027ea817e68f217488c",
					"transactions": []map[string]interface{}{
						testTransactionResponse,
					},
				},
			},
		}),
	}

	response, err := tclient.GetTransactions(&walletd.ReqGetTransactions{
		Addresses: []string{testAddress, testAddress},
		PaymentID: "test_payment_id",
		BlockHash: blockHash,
	})

	assert.Nil(err)

	item := response.Items[0]
	transaction := item.Transacations[0]

	assert.Equal("01bd06ca731914f27e143bbb902ce0bc05bff13d76faa027ea817e68f217488c", item.BlockHash)

	assert.Equal(-70368475742208, transaction.Fee)
	assert.Equal("0127cea59bfadc49aa02ed4a225936671e55607b5241621abca2a5e14405906dbb", transaction.Extra)
	assert.Equal(1446029698, transaction.Timestamp)
	assert.Equal(1, transaction.BlockIndex)
	assert.Equal(0, transaction.State)
	assert.Equal("06ec210a8359f253f8b2160a0d6040cf89f2a05a553aaa577b7f508ee5d831f9", transaction.TransactionHash)
	assert.Equal(70368475742208, transaction.Amount)
	assert.Equal(11, transaction.UnlockTime)
	assert.Equal("", transaction.PaymentID)
	assert.Equal(true, transaction.IsBase)

	assert.Equal(70368475742208, transaction.Transfers[0].Amount)
	assert.Equal(0, transaction.Transfers[0].Type)
	assert.Equal("KfXkT5VmdqmA7bWqSH37p87hSXBdTpTogN4mGHPARUSJaLse6jbXaVbVkLs3DwcmuD88xfu835Zvh6qBPCUXw6CHK8koDCt", transaction.Transfers[0].Address)
}

func TestGetUnconfirmedTransactionHashes(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "getUnconfirmedTransactionHashes",
		result: toJSONResult(map[string]interface{}{
			"transactionHashes": []string{
				"test_trans1",
				"test_trans2",
			},
		}),
	}

	response, err := tclient.GetUnconfirmedTransactionHashes(nil)

	assert.Nil(t, err)
	assert.Equal(t, "test_trans1", response.TransactionHashes[0])
	assert.Equal(t, "test_trans2", response.TransactionHashes[1])

	reqAssertion = requestAssertion{
		method: "getUnconfirmedTransactionHashes",
		params: map[string]interface{}{
			"addresses": []interface{}{testAddress, testAddress},
		},
		result: toJSONResult(map[string]interface{}{
			"transactionHashes": []interface{}{
				"test_trans1",
				"test_trans2",
			},
		}),
	}

	response, err = tclient.GetUnconfirmedTransactionHashes(&walletd.ReqGetUnconfirmedTransactionHashes{
		Addresses: []string{testAddress, testAddress},
	})

	assert.Nil(t, err)
	assert.Equal(t, "test_trans1", response.TransactionHashes[0])
	assert.Equal(t, "test_trans2", response.TransactionHashes[1])
}

func TestGetTransaction(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "getTransaction",
		params: map[string]interface{}{
			"transactionHash": "trans_hash1",
		},
		result: toJSONResult(map[string]interface{}{
			"transaction": testTransactionResponse,
		}),
	}

	response, err := tclient.GetTransaction("trans_hash1")
	transaction := response.Transaction

	assert.Nil(t, err)

	fmt.Printf("Transaction: %v\n", response)

	assert.Equal(t, -70368475742208, transaction.Fee)
	assert.Equal(t, "0127cea59bfadc49aa02ed4a225936671e55607b5241621abca2a5e14405906dbb", transaction.Extra)
	assert.Equal(t, 1446029698, transaction.Timestamp)
	assert.Equal(t, 1, transaction.BlockIndex)
	assert.Equal(t, 0, transaction.State)
	assert.Equal(t, "06ec210a8359f253f8b2160a0d6040cf89f2a05a553aaa577b7f508ee5d831f9", transaction.TransactionHash)
	assert.Equal(t, 70368475742208, transaction.Amount)
	assert.Equal(t, 11, transaction.UnlockTime)
	assert.Equal(t, "", transaction.PaymentID)
	assert.Equal(t, true, transaction.IsBase)
}

func TestGetTransactionSecretKey(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "getTransactionSecretKey",
		params: map[string]interface{}{
			"transaction": "trans_hash1",
		},
		result: toJSONResult(map[string]interface{}{
			"transactionSecretKey": "test_secret_key",
		}),
	}

	response, err := tclient.GetTransactionSecretKey("trans_hash1")

	assert.Nil(t, err)
	assert.Equal(t, "test_secret_key", response.TransactionSecretKey)
}

func TestGetTransactionProof(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "getTransactionProof",
		params: map[string]interface{}{
			"transactionHash":      "trans_hash1",
			"destinationAddress":   testAddress,
			"transactionSecretKey": "test_secret_key",
		},
		result: toJSONResult(map[string]interface{}{
			"transactionProof": "trans_proof",
		}),
	}

	response, err := tclient.GetTransactionProof(&walletd.ReqGetTransactionProof{
		TransactionHash:      "trans_hash1",
		DestinationAddress:   testAddress,
		TransactionSecretKey: "test_secret_key",
	})

	assert.Nil(t, err)
	assert.Equal(t, "trans_proof", response.TransactionProof)
}

func TestSendTransaction(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "sendTransaction",
		params: map[string]interface{}{
			"transfers": []interface{}{
				map[string]interface{}{
					"address": testAddress,
					"amount":  float64(1000000000000),
				},
			},
			"fee":       float64(10000000000),
			"anonymity": float64(3),
		},
		result: toJSONResult(map[string]interface{}{
			"transactionHash": "test_trans1",
		}),
	}

	response, err := tclient.SendTransaction(&walletd.ReqSendTransaction{
		Transfers: []struct {
			Address string "json:\"address\""
			Amount  uint64 "json:\"amount\""
		}{
			{
				Address: testAddress,
				Amount:  1000000000000,
			},
		},
		Fee:       10000000000,
		Anonymity: 3,
	})

	assert.Nil(t, err)
	assert.Equal(t, "test_trans1", response.TransactionHash)
}

func TestCreateDelayedTransaction(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "createDelayedTransaction",
		params: map[string]interface{}{
			"transfers": []interface{}{
				map[string]interface{}{
					"address": testAddress,
					"amount":  float64(1000000000000),
				},
			},
			"fee":       float64(10000000000),
			"anonymity": float64(3),
		},
		result: toJSONResult(map[string]interface{}{
			"transactionHash": "test_trans1",
		}),
	}

	response, err := tclient.CreateDelayedTransaction(&walletd.ReqCreateDelayedTransaction{
		Transfers: []struct {
			Address string "json:\"address\""
			Amount  uint64 "json:\"amount\""
		}{
			{
				Address: testAddress,
				Amount:  1000000000000,
			},
		},
		Fee:       10000000000,
		Anonymity: 3,
	})

	assert.Nil(t, err)
	assert.Equal(t, "test_trans1", response.TransactionHash)
}

func TestGetDelayedTransactionHashes(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "getDelayedTransactionHashes",
		result: toJSONResult(map[string]interface{}{
			"transactionHashes": []string{"test_hash1", "test_hash2"},
		}),
	}

	response, err := tclient.GetDelayedTransactionHashes()

	assert.Nil(t, err)
	assert.Equal(t, "test_hash1", response.TransactionHashes[0])
	assert.Equal(t, "test_hash2", response.TransactionHashes[1])
}

func TestDeleteDelayedTransaction(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "deleteDelayedTransaction",
		params: map[string]interface{}{
			"transactionHash": "test_hash",
		},
		result: "{}",
	}

	err := tclient.DeleteDelayedTransaction("test_hash")

	assert.Nil(t, err)
}

func TestSendDelayedTransaction(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "sendDelayedTransaction",
		params: map[string]interface{}{
			"transactionHash": "test_hash",
		},
		result: "{}",
	}

	err := tclient.SendDelayedTransaction("test_hash")

	assert.Nil(t, err)
}

func TestGetViewKey(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "getViewKey",
		result: toJSONResult(map[string]interface{}{
			"viewSecretKey": "test_secret_key",
		}),
	}

	response, err := tclient.GetViewKey()

	assert.Nil(t, err)
	assert.Equal(t, "test_secret_key", response.ViewSecretKey)
}

func TestGetMnemonicSeed(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "getMnemonicSeed",
		params: map[string]interface{}{
			"address": testAddress,
		},
		result: toJSONResult(map[string]interface{}{
			"mnemonicSeed": "test_seed",
		}),
	}

	response, err := tclient.GetMnemonicSeed(testAddress)

	assert.Nil(t, err)
	assert.Equal(t, "test_seed", response.MnemonicSeed)
}

func TestGetStatus(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "getStatus",
		result: toJSONResult(map[string]interface{}{
			"blockCount":            float64(256739),
			"knownBlockCount":       float64(256739),
			"localDaemonBlockCount": float64(256739),
			"lastBlockHash":         "1ae2b3fd7351a84c775e3088869efcc8afe6424cc2bc0ba9bd448c542061099b",
			"peerCount":             float64(5),
			"minimalFee":            float64(15645551732),
			"version":               "1.7.3.923-eb91ef4d9",
		}),
	}

	response, err := tclient.GetStatus()

	assert.Nil(t, err)

	assert.Equal(t, uint32(256739), response.BlockCount)
	assert.Equal(t, uint32(256739), response.KnownBlockCount)
	assert.Equal(t, uint32(256739), response.LocalDaemonBlockCount)
	assert.Equal(t, "1ae2b3fd7351a84c775e3088869efcc8afe6424cc2bc0ba9bd448c542061099b", response.LastBlockHash)
	assert.Equal(t, uint32(5), response.PeerCount)
	assert.Equal(t, uint64(15645551732), response.MinimalFee)
	assert.Equal(t, "1.7.3.923-eb91ef4d9", response.Version)
}

func TestGetAddresses(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "getAddresses",
		result: toJSONResult(map[string]interface{}{
			"addresses": []string{"test_address1", "test_address2"},
		}),
	}

	response, err := tclient.GetAddresses()

	assert.Nil(t, err)

	assert.Equal(t, "test_address1", response.Addresses[0])
	assert.Equal(t, "test_address2", response.Addresses[1])
}

func TestGetAddressesCount(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "getAddressesCount",
		result: toJSONResult(map[string]interface{}{
			"addressesCount": 3,
		}),
	}

	response, err := tclient.GetAddressesCount()

	assert.Nil(t, err)

	assert.Equal(t, uint32(3), response.AddressesCount)
}

func TestSendFusionTransaction(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "sendFusionTransaction",
		params: map[string]interface{}{
			"threshold":          float64(1000000),
			"anonymity":          float64(4),
			"addresses":          []interface{}{testAddress, testAddress},
			"destinationAddress": testAddress,
		},
		result: toJSONResult(map[string]interface{}{
			"transactionHash": "tran_hash",
		}),
	}

	response, err := tclient.SendFusionTransaction(&walletd.ReqSendFusionTransaction{
		Threshold:          1000000,
		Anonymity:          4,
		Addresses:          []string{testAddress, testAddress},
		DestinationAddress: testAddress,
	})

	assert.Nil(t, err)
	assert.Equal(t, "tran_hash", response.TransactionHash)
}

func TestEstimateFusion(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "estimateFusion",
		params: map[string]interface{}{
			"threshold": float64(1000000),
			"addresses": []interface{}{testAddress, testAddress},
		},
		result: toJSONResult(map[string]interface{}{
			"totalOutputCount": 1000,
			"fusionReadyCount": 50,
		}),
	}

	response, err := tclient.EstimateFusion(&walletd.ReqEstimateFusion{
		Threshold: 1000000,
		Addresses: []string{testAddress, testAddress},
	})

	assert.Nil(t, err)
	assert.Equal(t, 1000, response.TotalOutputCount)
	assert.Equal(t, 50, response.FusionReadyCount)
}

func TestValidateAddress(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "validateAddress",
		params: map[string]interface{}{
			"address": testAddress,
		},
		result: toJSONResult(map[string]interface{}{
			"address":        testAddress,
			"isValid":        true,
			"spendPublicKey": "spend_public_key",
			"viewPublicKey":  "view_public_key",
		}),
	}

	response, err := tclient.ValidateAddress(testAddress)

	assert.Nil(t, err)
	assert.Equal(t, testAddress, response.Address)
	assert.Equal(t, true, response.IsValid)
	assert.Equal(t, "spend_public_key", response.SpendPublicKey)
	assert.Equal(t, "view_public_key", response.ViewPublicKey)
}

func TestGetReserveProof(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "reserveProof",
		params: map[string]interface{}{
			"address": testAddress,
			"message": "test_message",
			"amount":  float64(1000),
		},
		result: toJSONResult(map[string]string{
			"reserveProof": "test_reserve_proof",
		}),
	}

	response, err := tclient.GetReserveProof(&walletd.ReqGetReserveProof{
		Address: testAddress,
		Message: "test_message",
		Amount:  1000,
	})

	assert.Nil(t, err)
	assert.Equal(t, "test_reserve_proof", response.ReserveProof)
}

func TestSignMessage(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "signMessage",
		params: map[string]interface{}{
			"address": testAddress,
			"message": "test_message",
		},
		result: toJSONResult(map[string]interface{}{
			"address":   testAddress,
			"signature": "test_signature",
		}),
	}

	response, err := tclient.SignMessage(&walletd.ReqSignMessage{
		Address: testAddress,
		Message: "test_message",
	})

	assert.Nil(t, err)
	assert.Equal(t, testAddress, response.Adddress)
	assert.Equal(t, "test_signature", response.Signature)
}

func TestVerifyMessage(t *testing.T) {
	tclient := buildTestClient(t)

	reqAssertion = requestAssertion{
		method: "verifyMessage",
		params: map[string]interface{}{
			"address":   testAddress,
			"message":   "test_message",
			"signature": "test_signature",
		},
		result: toJSONResult(map[string]interface{}{
			"isValid": true,
		}),
	}

	response, err := tclient.VerifyMessage(&walletd.ReqVerifyMessage{
		Address:   testAddress,
		Message:   "test_message",
		Signature: "test_signature",
	})

	assert.Nil(t, err)
	assert.True(t, response.IsValid)
}
