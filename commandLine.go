package main

import "fmt"

func (cli *CLI) AddBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Printf("add block success \n")
}

func (cli *CLI) PrintBlockChain() {
	bc := cli.bc
	// 调用迭代器 返回每一个数据
	it := bc.NewIterator()
	for {
		block := it.Next()
		fmt.Printf("===========当前区块高度：  ==============\n")
		fmt.Printf("前区块哈希： %x\n", block.PrevHash)
		fmt.Printf("当前区块哈希： %x\n", block.Hash)
		fmt.Printf("区块数据： %s\n", block.Data)
		if len(block.PrevHash) == 0 {
			fmt.Printf("over")
			break
		}
	}
}
