package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"log"
)

const reward = 50

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
	//Sig string
	// 真正的数字签名，由r s 拼成的[]byte
	Signature []byte
	// 这里的PubKey不存储原始的公钥，而是存储X和Y拼接的字符串
	PubKey []byte
}
type TXOutput struct {
	// 转账金额
	Value float64
	// 锁定脚本，用地址模拟
	//PubKeyHash string
	// 收款方的公钥的哈希，是哈希
	PubKeyHash []byte
}

// 由于现在存储的字段是地址的公钥哈希，所以无法直接创建TXOutput
// 为了能够得到公钥哈希，我们需要处理一下，写一个Lock函数
func (output *TXOutput) Lock(address string) {
	// 1. 解码
	// 2. 截取出公钥哈希：去除version（1字节）去除校验码（4字节）
	addressByte := base58.Decode(address) //25字节
	len := len(addressByte)
	pubKeyHash := addressByte[1 : len-4]
	// 真正的锁定动作
	output.PubKeyHash = pubKeyHash

}

// 给TXOutput一个创建你的方法，否则无法调用Lock
func NewTXOutput(value float64, address string) *TXOutput {
	output := TXOutput{
		Value: value,
	}
	output.Lock(address)
	return &output
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
	// 矿工挖矿时无需指定签名，因此PubKey可以自由填写。一般填写矿工（矿池）的名字
	// 签名先填写为空，后面创建完整交易后最后做一次签名即可
	input := TXInput{[]byte{}, -1, nil, []byte(data)}
	//output := TXOutput{reward, address}
	output := NewTXOutput(reward, address)
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{*output}}
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
	// 1. 交易数字签名-》私钥-》打开钱包-》NewWallets()
	ws := NewWallets()
	// 2. 根据地址 返回自己的wallet
	wallet := ws.WalletsMap[from]
	if wallet == nil {
		fmt.Printf("没有找到该地址的钱包，交易创建失败\n")
		return nil
	}
	// 3. 得到对应的公钥 私钥
	pubKey := wallet.PubKey
	//privateKey:=wallet.Private // 稍后再用

	// 输入：2*地址，金额

	// 传递公钥的哈希，而不是传递地址
	// func HashPubKey(data []byte) []byte {
	pubKeyHash := HashPubKey(pubKey)
	// 找到合理的返回
	// func (bc *BlockChain) FindNeedUTXOs(from string, amount float64) (map[string][]uint64, float64) {
	utxos, resValue := bc.FindNeedUTXOs(pubKeyHash, amount)
	if resValue < amount {
		fmt.Printf("余额不足 失败\n")
		return nil
	}
	var inputs []TXInput
	var outputs []TXOutput
	// 2. 创建交易输入，将utxo逐一转成input
	for id, indexArray := range utxos {
		for _, i := range indexArray {
			// 签名需要私钥，私钥在本地的钱包里，需要拿出来
			input := TXInput{[]byte(id), int64(i), nil, pubKey}
			inputs = append(inputs, input)
		}
	}
	//output := TXOutput{amount, to}
	output := NewTXOutput(amount, to)
	outputs = append(outputs, *output)
	if resValue > amount {
		// 找零
		output = NewTXOutput(resValue-amount, from)
		outputs = append(outputs, *output)
	}
	tx := Transaction{[]byte{}, inputs, outputs}
	tx.SetHash()
	return &tx
}
