package main

import (
	"fmt"
	"os"
)

// 用来接收命令行参数并且控制区块链操作的文件

type CLI struct {
	bc *BlockChain
}

const Usage = `
	addBlock --data DATA "add data to blockchain"
	printChain           "print all blockchain data"
`

// 接受参数的动作放到一个函数中
func (cli *CLI) Run() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Invalid command")
		fmt.Printf(Usage)
	}
	cmd := args[1]
	switch cmd {
	case "addBlock":
		// 添加区块
		fmt.Printf("add block")

	case "printChain":
		// 打印区块
		fmt.Printf("print block")
	default:
		fmt.Printf("Invalid command")
		fmt.Printf(Usage)
	}

	// 分析命令
	// 1. 添加区块 2. 打印区块

	// 执行相应动作
}
