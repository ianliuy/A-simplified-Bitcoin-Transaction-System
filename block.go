package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
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
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
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
	//block.SetHash()
	// 创建一个pow对象
	pow := NewProofOfWork(&block)
	// 查找随机数 不停进行哈希运算
	hash, nonce := pow.Run()

	// 根据挖矿结果 不断对区块数据进行补充
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

func (block *Block) toByte() []byte {
	//TODO
	return []byte{}
}

// 3. 生成哈希
func (block *Block) SetHash() {
	var blockInfo []byte
	// 1. 拼装数据
	/*blockInfo = append(block.PrevHash, block.Data...)
	blockInfo = append(blockInfo, Uint64ToByte(block.Version)...)
	blockInfo = append(blockInfo, block.PrevHash...)
	blockInfo = append(blockInfo, block.MerkelRoot...)
	blockInfo = append(blockInfo, Uint64ToByte(block.TimeStamp)...)
	blockInfo = append(blockInfo, Uint64ToByte(block.Difficulty)...)
	blockInfo = append(blockInfo, Uint64ToByte(block.Nonce)...)
	blockInfo = append(blockInfo, block.Data...)*/
	tmp := [][]byte{
		Uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(block.Nonce),
		block.Data,
	}

	blockInfo = bytes.Join(tmp, []byte{})
	// 2. sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
