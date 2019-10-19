package main

import (
	"crypto/sha256"
)

// 0. 定义结构
type Block struct {
	// 1. 前区块哈希
	PrevHash []byte
	// 2. 当前区块哈希
	Hash []byte
	// 3. 数据
	Data []byte
}

// 2. 创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		PrevHash: prevBlockHash,
		Hash:     []byte{}, //先空 //TODO
		Data:     []byte(data),
	}
	block.SetHash()
	return &block
}

// 3. 生成哈希
func (block *Block) SetHash() {
	// 1. 拼装数据
	blockInfo := append(block.PrevHash, block.Data...)
	// 2. sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
