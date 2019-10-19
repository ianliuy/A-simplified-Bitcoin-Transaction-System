package main

import (
	"crypto/sha256"
	"time"
)

// 0. 定义结构
type Block struct {
	// 版本号
	Version uint64
	// 1. 前区块哈希
	PrevHash []byte
	// Merkelroot(梅克尔根，这就是一个哈希值)
	MerkelRoot []byte
	// 时间戳
	TimeStamp uint64
	// 难度值
	Difficulty uint64
	// 随机数 挖矿要找的数
	Nonce uint64

	// 2. 当前区块哈希 //正常比特币区块中没有当前区块的哈希，为了简化才做这些
	Hash []byte
	// 3. 数据
	Data []byte
}

// 1. 补充区块其他字段
// 2. 更新计算哈希函数
// 3. 优化代码

// 实现一个辅助函数，将uint转换成[]byte
func Uint64ToByte(num uint64) []byte {
	//TODO
	return []byte{}
}

// 2. 创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,        //随便写的无效值
		Nonce:      0,        // 无效值
		Hash:       []byte{}, //先空 //TODO
		Data:       []byte(data),
	}
	block.SetHash()
	return &block
}

// 3. 生成哈希
func (block *Block) SetHash() {
	var blockInfo []byte
	// 1. 拼装数据
	blockInfo = append(block.PrevHash, block.Data...)
	blockInfo = append(blockInfo, Uint64ToByte(block.Version)...)
	blockInfo = append(blockInfo, block.PrevHash...)
	blockInfo = append(blockInfo, block.MerkelRoot...)
	blockInfo = append(blockInfo, Uint64ToByte(block.TimeStamp)...)
	blockInfo = append(blockInfo, Uint64ToByte(block.Difficulty)...)
	blockInfo = append(blockInfo, Uint64ToByte(block.Nonce)...)
	blockInfo = append(blockInfo, block.Data...)
	// 2. sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
