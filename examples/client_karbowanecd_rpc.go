package main

import (
	"fmt"

	"github.com/r3volut1oner/go-karbo/client/karbowanecd"
)

func main() {

	client := karbowanecd.NewClient("http://localhost:32348/json_rpc")

	address := "KeCp19gGMUPg1wVWeagX3AUcK1ZmCxtqsEVGuNzkAgnyeftV4whKUBQShUmPPqiBZmeeX1TX2en6qgbzFtYUy1xG27gHTst"

	height, err := client.GetBlockCount()
	if err != nil {
		panic(err)
	}

	fmt.Println("Block count:", height)

	hash, err := client.GetBlockHash(height - 200)
	if err != nil {
		panic(err)
	}

	fmt.Println("Block hash:", hash)

	blockTemplate, err := client.GetBlockTemplate(2, address)
	if err != nil {
		panic(err)
	}

	fmt.Println("Block template hashing:", blockTemplate.BlockHashingBlob)
	fmt.Println()

	header, err := client.GetBlockHeaderByHash(hash)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Block header by hash: %+v\n", header)
	fmt.Println("Block header timestamp:", header.Timestamp)
	fmt.Println()

	header, err = client.GetBlockHeaderByHeight(height - 2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Block header by height: %+v\n", header)
	fmt.Println()

	timestamp, err := client.GetBlockTimestamp(height - 100)
	if err != nil {
		panic(err)
	}

	fmt.Println("Block timestamp:", timestamp)
	fmt.Println()

	// block, err := client.GetBlockByHeight(height - 2)
	// if err != nil {
	// 	panic(err)
	// }

	// printBlock(block)
	// fmt.Println()

	// hash = "fa557ff5eeab183ef38b56c0d29db2729ac7c3a32fc1351fd0347b69f6127aef"

	block, err := client.GetBlockByHash(hash)
	if err != nil {
		panic(err)
	}

	printBlock(&block)

	fmt.Println()

	blocks, err := client.GetBlocksByHeights([]int{height - 3, height - 2, height - 1})
	if err != nil {
		panic(err)
	}

	fmt.Println("Blocks by height number:", len(blocks))

	blocks, err = client.GetBlocksByHashes([]string{blocks[0].Hash, blocks[1].Hash, blocks[2].Hash})
	if err != nil {
		panic(err)
	}

	fmt.Println("Blocks hash number:", len(blocks))

	tHashes, count, err := client.GetBlocksHashesByTimestamps(timestamp, timestamp+7200, 1)
	if err != nil {
		panic(err)
	}

	fmt.Println("Hashes by timestamps:", len(tHashes), "Count:", count)

	blocksList, err := client.GetBlocksList(height-20, 0)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Blocks list: %+v\n", blocksList)
	fmt.Println()

	altBlocks, err := client.GetAltBlocksList()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Alt. blocks list: %+v\n", altBlocks)
	fmt.Println()

	lastBlockHeader, err := client.GetLastBlockHeader()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Last block header: %+v\n", lastBlockHeader)
	fmt.Println()

	thash := "a1d42c5a01950ae5f89cffb10760f890c3f25c8d209c82317ecde1c712d09083"

	transaction, err := client.GetTransaction(thash)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transaction: %+v\n", transaction)
	fmt.Println()

	mempool, err := client.GetTransactionsPool()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Mempool: %+v\n", mempool)
	fmt.Println()

	paymentID := "69402fedb764b375c5de9202a6fd0dd4bc2319e217a8a9245ffb576bac473f2f"
	shortTransactions, err := client.GetTransactionsByPaymentID(paymentID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Payment id \"%s\" transactions length: %+v\n", paymentID, len(shortTransactions))
	fmt.Println()

	tHashes, err = client.GetTransactionHashesByPaymentID(paymentID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Payment id \"%s\" transactions hashes length: %+v\n", paymentID, len(tHashes))
	fmt.Println()

	transactions, err := client.GetTransactionsByHashes([]string{thash})
	if err != nil {
		panic(err)
	}

	fmt.Println("Transactions by hashes length:", len(transactions))
	fmt.Println()

	currid, err := client.GetCurrencyID()
	if err != nil {
		panic(err)
	}

	fmt.Println("Currency ID:", currid)
	fmt.Println()

	stats, err := client.GetStatsByHeights([]int{height - 2, height - 3})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Stats: %+v\n", stats)
	fmt.Println()

	stats, err = client.GetStatsByHeights([]int{height - 2, height - 5})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Stats in range: %+v\n", len(stats))
	fmt.Println()

	validation, err := client.ValidateAddress(address)

	fmt.Printf("Address validation: %+v\n", validation)
	fmt.Println()
}

func printBlock(block *karbowanecd.Block) {
	fmt.Println("Block Hash:", block.Hash)

	for tidx, t := range block.Transactions {
		fmt.Println("Tranaction hash:", tidx, t.Hash)
		fmt.Println()

		for iidx, i := range t.Inputs {
			fmt.Println("  Input type:", iidx, i.Type)
			fmt.Printf("  Input data: %+v\n", i.Data)
		}

		fmt.Println()

		for oidx, o := range t.Outputs {
			fmt.Println("  Output global index:", oidx, o.GlobalIndex)
			fmt.Printf("  Output data: %+v\n", o.Output)
		}

		fmt.Println()

		fmt.Printf("  Extra: %+v\n", t.Extra)

		fmt.Println()
	}

	fmt.Println()
	fmt.Printf("Block by hash: %+v\n", block)

}
