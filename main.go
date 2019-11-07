package main

// import "fmt"

// 6. 重构代码

func main() {
	bc := NewBlockChain("19g9h5Xh24Vb92TCwmj1yxL88KwiZ5ytdZ")
	cli := CLI{bc}
	cli.Run()

	/*	bc.AddBlock("first transaction")
		bc.AddBlock("second transaction")

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
		}*/
}
