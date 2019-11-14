package main

import (
	"fmt"
	"time"
)

func (cli *CLI) PrintBlockChain() {
	bc := cli.bc
	// 调用迭代器 返回每一个数据
	blockHeight := bc.GetBlockHeight() // get the height of the Blockchain
	it := bc.NewIterator()
	for {
		block := it.Next()
		fmt.Printf("===========Height: %v ==============\n", blockHeight)
		fmt.Printf("Version: %d\n", block.Version)
		fmt.Printf("Prev Block Hash: %x\n", block.PrevHash)
		fmt.Printf("Merkle root: %x\n", block.MerkelRoot)
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("Timestamp: %s\n", timeFormat)
		fmt.Printf("Difficulty: %d\n", block.Difficulty)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Info: %s\n", block.Transactions[0].TXInputs[0].PubKey)
		if len(block.PrevHash) == 0 {
			fmt.Printf("over")
			break
		}
		blockHeight--
	}
}

func (cli *CLI) GetBalance(address string) {
	pubKeyHash := GetPubeyFromAddress(address)
	utxos := cli.bc.FindUTXOs(pubKeyHash)
	total := 0.0
	for _, utxos := range utxos {
		total += utxos.Value
	}
	fmt.Printf("\"%s\"'s balance is: %f\n", address, total)
}

func (cli *CLI) Send(from, to string, amount float64, miner, data string) {
	// 1. creat a Coinbase transaction
	coinbase := NewCoinbaseTX(miner, data)
	// 2. creat a normal transaction
	tx := NewTransaction(from, to, amount, cli.bc)
	// 3. add block
	cli.bc.AddBlock([]*Transaction{coinbase, tx})
	//cli.GetBalance("1NLPgzsGC79JLFwmyoUpvhBrWRzES7vuSV")
	fmt.Printf("transact success\n")
}

func (cli *CLI) NewWallet() {
	ws := NewWallets()
	address := ws.CreatWallet()
	fmt.Printf("address:%v\n", address)
}

func (cli *CLI) listAddresses() {
	ws := NewWallets()
	addresses := ws.ListAllAddresses()
	for i, address := range addresses {
		fmt.Printf("address %v: %v           ", i, address)
		cli.GetBalance(address)
	}
}
