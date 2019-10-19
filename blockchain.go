package main

// 4. 引入区块
type BlockChain struct {
	// 定义一个区块链数组
	blocks []*Block
}

// 5. 定义一个区块链
func NewBlockChain() *BlockChain {
	// 创建一个创世块并作为第一个区块添加到区块链中
	genesisBlock := GenesisBlock()
	return &BlockChain{
		blocks: []*Block{genesisBlock},
	}
}

// 创世块
func GenesisBlock() *Block {
	return NewBlock("创世块", []byte{})
}

// 5. 添加区块
func (bc *BlockChain) AddBlock(data string) {
	// 根据下标获取前区块哈希
	lastBlock := bc.blocks[len(bc.blocks)-1]
	prevHash := lastBlock.Hash
	// 1. 创建新区块
	block := NewBlock(data, prevHash)
	// 2. 添加到区块链数组中
	bc.blocks = append(bc.blocks, block)

}
