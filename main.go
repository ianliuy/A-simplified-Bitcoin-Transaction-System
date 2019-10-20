package main

// 6. 重构代码

func main() {
	bc := NewBlockChain()
	bc.AddBlock("first transaction")
	bc.AddBlock("second transaction")
	// 迭代器
	/*for i, block := range bc.blocks {
		fmt.Printf("===========当前区块高度： %d ==============\n", i)
		fmt.Printf("前区块哈希： %x\n", block.PrevHash)
		fmt.Printf("当前区块哈希： %x\n", block.Hash)
		fmt.Printf("区块数据： %s\n", block.Data)
	}*/
}
