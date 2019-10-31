package main

import (
	"fmt"
	"os"
	"strconv"
)

// 用来接收命令行参数并且控制区块链操作的文件

type CLI struct {
	bc *BlockChain
}

const Usage = `
	printChain                             "print all blockchain data"
    getBalance    --address  ADDRESS       "获取指定地址的余额"
    send FROM TO AMOUNT MINER DATA         "由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"
    newWallet                              "添加一个新的钱包"
    listAddresses                          "列举所有的钱包地址"
`

// 接受参数的动作放到一个函数中
func (cli *CLI) Run() {
	args := os.Args
	if len(args) < 2 {
		// fmt.Printf("Invalid command")
		// fmt.Printf(Usage)
	}
	// 2. 分析命令
	cmd := args[0]
	if len(args) > 1 {
		cmd = args[1]
	}
	switch cmd {
	case "printChain":
		// 打印区块
		fmt.Printf("print block\n")
		cli.PrintBlockChain()
	case "getBalance":
		fmt.Printf("获取余额\n")
		if len(args) == 4 && args[2] == "--address" {
			address := args[3]
			cli.GetBalance(address)
		}
	case "send":
		fmt.Printf("start sending...\n")
		if len(args) != 7 {
			fmt.Printf("参数个数错误\n")
			fmt.Printf(Usage)
			return
		}
		// send FROM TO AMOUNT MINER DATA         "由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"

		from := args[2]
		to := args[3]
		// func ParseFloat(s string, bitSize int) (float64, error) {
		amount, _ := strconv.ParseFloat(args[4], 64)
		miner := args[5]
		data := args[6]
		fmt.Printf("转账信息： %s %s %f %s %s\n", from, to, amount, miner, data)
		cli.Send(from, to, amount, miner, data)
	case "newWallet":
		fmt.Printf("creating new wallet.......\n")
		cli.NewWallet()
	case "listAddresses":
		fmt.Printf("列举所有钱包的地址\n")
		cli.listAddresses()
	default:
		// fmt.Printf("Invalid command")
		// fmt.Printf(Usage)
	}
	// 分析命令
	// 1. 添加区块 2. 打印区块
	// 执行相应动作
}
