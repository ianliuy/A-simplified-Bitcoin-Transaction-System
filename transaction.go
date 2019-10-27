package main

// 1. 定义交易结构
type Transaction struct {
	TXID      []byte     // 交易ID
	TXInputs  []TXInput  // 交易输入数组
	TXOutputs []TXOutput // 交易输出数组
}

// 定义交易输入
type TXInput struct {
	// 引用的交易ID
	TXid []byte
	// 引用的output的索引值
	Index int64
	// 解锁脚本，用地址来模拟
	Sig string
}

type TXOutput struct {
	// 转账金额
	value float64
	// 锁定脚本，用地址模拟
	PubKeyHash string
}

// 2. 提供创建交易的方法

// 3. 创建挖矿交易

// 4. 根据交易调整程序

//
