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
	addBlock --data DATA          "add data to blockchain"
	printChain                    "print all blockchain data"
`

// 接受参数的动作放到一个函数中
func (cli *CLI) Run() {
	// ./block printChain
	// ./block addBlock --data "HelloWorld"
	// 1. 得到所有的命令
	args := os.Args
	/*fmt.Printf("len(args)=%v\n", len(args))
	fmt.Printf("args[0]=%v\n", args[0])
	fmt.Printf("args[1]=%v\n", args[1])
	fmt.Printf("args[2]=%v\n", args[2])
	fmt.Printf("args[3]=%v\n", args[3])
	fmt.Printf("args[4]=%v\n", args[4])
	fmt.Println("args[1]:", args[1])*/
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
	case "addBlock":
		// 添加区块
		fmt.Printf("add block")

		// 确保命令有效
		if len(args) == 4 && args[2] == "--data" {
			// 获取命令的数据
			// 1. 获取数据
			data := args[3]
			cli.AddBlock(data)
		} else {
			fmt.Printf("errors occur")
			fmt.Printf(Usage)
		}

		// 2. 使用bc添加区块AddBlock

	case "printChain":
		// 打印区块
		fmt.Printf("print block\n")
		cli.PrintBlockChain()
	default:
		// fmt.Printf("Invalid command")
		// fmt.Printf(Usage)
	}
	// 分析命令
	// 1. 添加区块 2. 打印区块
	// 执行相应动作
}
