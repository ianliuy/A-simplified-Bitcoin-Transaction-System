package main

import (
	"fmt"
	"time"
)

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
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("时间戳： %s\n", timeFormat)
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
	//fmt.Printf("func (cli *CLI) Send(from:%s, to:%s string, amount: %f float64, miner: %s, data: %s string) {\n", from, to, amount, miner, data)
	//fmt.Printf("经过func (cli *CLI) Send，张三的余额为：\n")
	//cli.GetBalance("张三")
	// 1. 创建挖矿交易
	coinbase := NewCoinbaseTX(miner, data)
	//fmt.Printf("coinbase := NewCoinbaseTX(miner:%s, data:%s)\n", miner, data)
	//fmt.Printf("经过NewCoinbaseTX，张三的余额为：\n")
	//cli.GetBalance("张三")
	// 2. 创建一个普通交易
	tx := NewTransaction(from, to, amount, cli.bc)
	//fmt.Printf("tx := NewTransaction(from:%s, to:%s, amount:%s, cli.bc)\n", from, to, amount)
	//fmt.Printf("经过NewTransaction，张三的余额为：\n")
	//cli.GetBalance("张三")
	// 3. 添加区块
	cli.bc.AddBlock([]*Transaction{coinbase, tx})
	//fmt.Printf("cli.bc.AddBlock([]*Transaction{coinbase, tx}), coinbase.TXID:%v\n", string(coinbase.TXID))
	fmt.Printf("经过cli.bc.AddBlock,")
	cli.GetBalance("张三")
	fmt.Printf("transact success\n")
}

func (cli *CLI) NewWallet() {
	ws := NewWallets()
	address := ws.CreatWallet()
	fmt.Printf("address:%v", address)
	// func NewWallet() *Wallet {
	//wallet := NewWallet()
	//address := wallet.NewAddress()
	/*	ws := NewWallets()
		for address, _ := range ws.WalletsMap {
			fmt.Printf("地址：%v\n", address)
		}*/
	//fmt.Printf("私钥：%v\n", wallet.Private)
	//fmt.Printf("公钥：%v\n", wallet.PubKey)
	//fmt.Printf("地址：%v\n", address)
}

func (cli *CLI) listAddresses() {
	ws := NewWallets()
	addresses := ws.ListAllAddresses()
	for i, address := range addresses {
		fmt.Printf("address %v: %v\n", i, address)
	}
}
