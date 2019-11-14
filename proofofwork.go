package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strconv"
)

// 定义一个工作量证明的结构
type ProofOfWork struct {
	// a. block
	block *Block
	// b. 目标值
	target *big.Int
}

// 2. 提供创建pow的函数
// NewProofofWork(参数)
func NewProofOfWork(block *Block, difficulty uint64) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}
	// 指定难度值
	// targetStr := GetTargetStr()
	targetStr := strconv.Itoa(int(difficulty))
	for {
		if len(targetStr) == 64 {
			break
		}
		targetStr = string('0') + targetStr
	}
	targetStr = "0010000000000000000000000000000000000000000000000000000000000000"
	// 引入辅助变量。str->big.int
	temInt := big.Int{}
	// 16进制格式
	temInt.SetString(targetStr, 16)

	pow.target = &temInt
	return &pow
}

// 提供不断计算hash的函数
// Run
func (pow *ProofOfWork) Run() ([]byte, uint64) {

	var nonce uint64
	block := pow.block
	var hash [32]byte
	fmt.Printf("start mining: ")
	for {
		// 1. 拼装数据（区块、随机数）
		tmp := [][]byte{
			Uint64ToByte(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			Uint64ToByte(block.TimeStamp),
			Uint64ToByte(block.Difficulty),
			Uint64ToByte(nonce),
			// 只对区块头做哈希值，区块提通过Merkelroot产生影响
			//block.Data,

		}
		blockInfo := bytes.Join(tmp, []byte{})

		// 2. 哈希运算
		hash = sha256.Sum256(blockInfo)
		// 3. 跟数target比较
		temInt := big.Int{}
		// 将我们得到的hash数组转化成一个big.Int
		temInt.SetBytes(hash[:])

		//比较当前的哈希与目标哈希值。小于就i找到了，没小于i就继续找
		if temInt.Cmp(pow.target) == -1 {
			fmt.Printf("mining success!! hash: %x, nonce: %d\n", hash, nonce)
			//break
			return hash[:], nonce
		} else {
			//	没找到 继续找
			nonce++
		}
		//	找到了：退出返回
		//  没找到：继续找，随机数加一
	}

	//return hash[:], nonce
}
