package main

import "math/big"

// 定义一个工作量证明的结构
type ProofOfWork struct {
	// a. block
	block *Block
	// b. 目标值
	target *big.Int
}

// 2. 提供创建pow的函数
// NewProofofWork(参数)
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}
	// 指定难度值
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"
	// 引入辅助变量。str->big.int
	temInt := big.Int{}
	// 16进制格式
	temInt.SetString(targetStr, 16)

	pow.target = &temInt
	return &pow
}

// 提供不断计算hash的函数

// Run

// IsValid
