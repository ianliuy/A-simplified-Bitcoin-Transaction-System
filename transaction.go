package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const reward = 12.5

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
	Value float64
	// 锁定脚本，用地址模拟
	PubKeyHash string
}

// 设置交易ID

func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}

// 实现一个函数，是否是挖矿交易
func (tx *Transaction) IsCoinbase() bool {
	// 交易input只有一个
	/*if len(tx.TXInputs) == 1 {
		input := tx.TXInputs[0]
		// 交易id为空
		// 交易index是-1
		if !bytes.Equal(tx.TXInputs[0].TXid, []byte{}) || input.Index != -1 {
			return false
		}
	}
	return true*/
	if len(tx.TXInputs) == 1 && len(tx.TXInputs[0].TXid) == 0 && tx.TXInputs[0].Index == -1 {
		return true
	}
	return false
}

// 2. 提供创建交易的方法(挖矿交易 coinbase)
func NewCoinbaseTX(address string, data string) *Transaction {
	// 挖矿交易的特点：只有一个input
	// 无需引用交易id
	// 无需引用index
	// 矿工挖矿时无需指定签名，因此sig可以自由填写。一般填写矿工（矿池）的名字
	input := TXInput{[]byte{}, -1, data}
	output := TXOutput{reward, address}
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{output}}
	tx.SetHash()
	return &tx
	// Transaction{}
}

// 创建普通转账交易
// 1. 找到最合理的UTXO集合（使用map标注，map[string][]uint64）
// 2. 将utxo逐一转成input
// 3. 创建outputs
// 4. 有零钱的话，找零
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	// 输入：2*地址，金额

	// 找到合理的返回
	// func (bc *BlockChain) FindNeedUTXOs(from string, amount float64) (map[string][]uint64, float64) {
	utxos, resValue := bc.FindNeedUTXOs(from, amount)
	if resValue < amount {
		fmt.Printf("余额不足 失败")
		return nil
	}
	var inputs []TXInput
	var outputs []TXOutput
	// 2. 创建交易输入，将utxo逐一转成input
	for id, indexArray := range utxos {
		for _, i := range indexArray {
			input := TXInput{[]byte(id), int64(i), from}
			inputs = append(inputs, input)
		}
	}
	output := TXOutput{amount, to}
	outputs = append(outputs, output)
	if resValue > amount {
		// 找零
		outputs = append(outputs, TXOutput{resValue - amount, from})
	}
	tx := Transaction{[]byte{}, inputs, outputs}
	tx.SetHash()
	return &tx
}

// 4. 根据交易调整程序

//
