package main

import (
	"fmt"
)

func (cli *CLI) AddBlock(data string) {
	//cli.bc.AddBlock(data)
	fmt.Printf("add block success \n")
}

func (cli *CLI) PrintBlockChain() {
	bc := cli.bc
	// 调用迭代器 返回每一个数据
	it := bc.NewIterator()
	for {
		block := it.Next()
		fmt.Printf("===========当前区块高度：  ==============\n")
		fmt.Printf("版本号： %d\n", block.Version)
		fmt.Printf("前区块哈希： %x\n", block.PrevHash)
		fmt.Printf("梅克尔根： %x\n", block.MerkelRoot)
		fmt.Printf("时间戳： %x\n", block.TimeStamp)
		fmt.Printf("难度值（随便写的）： %d\n", block.Difficulty)
		fmt.Printf("随机数： %d\n", block.Nonce)
		fmt.Printf("区块哈希： %x\n", block.Hash)
		fmt.Printf("区块数据： %s\n", block.Transactions[0].TXInputs[0].Sig)
		if len(block.PrevHash) == 0 {
			fmt.Printf("over")
			break
		}
	}
}

func (cli *CLI) GetBalance(address string) {
	utxos := cli.bc.FindUTXOs(address)
	total := 0.0
	for _, utxos := range utxos {
		total += utxos.Value
	}
	fmt.Printf("\"%s\"的余额为：%f\n", address, total)
}

func (cli *CLI) Send(from, to string, amount float64, miner, data string) {
	fmt.Printf("from: %s\n", from)
	fmt.Printf("to: %s\n", to)
	fmt.Printf("amount: %f\n", amount)
	fmt.Printf("miner: %s\n", miner)
	fmt.Printf("data: %s\n", data)
	// 具体的逻辑 TODO
	// 1. 创建挖矿交易
	coinbase := NewCoinbaseTX(miner, data)
	// 2. 创建一个普通交易
	tx := NewTransaction(from, to, amount, cli.bc)
	// 3. 添加区块
	cli.bc.AddBlock([]*Transaction{coinbase, tx})
	fmt.Printf("transact success")
}
