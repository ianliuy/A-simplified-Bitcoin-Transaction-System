package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"fmt"
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
	// 真实的交易数组
	Transactions []*Transaction
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
func NewBlock(txs []*Transaction, prevBlockHash []byte, bc *BlockChain) *Block {
	block := Block{
		Version:      00,
		PrevHash:     prevBlockHash,
		MerkelRoot:   []byte{},
		TimeStamp:    uint64(time.Now().Unix()),
		Difficulty:   0,
		Nonce:        0,
		Hash:         []byte{},
		Transactions: txs,
	}
	block.MerkelRoot = block.MakeMerkelRoot()
	//block.SetHash()

	heightNow := bc.GetBlockHeight() + 1
	timeStamp := block.TimeStamp
	difficulty := block.GetNowDifficulty(bc, heightNow, timeStamp)
	block.Difficulty = difficulty

	pow := NewProofOfWork(&block, difficulty)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	block.Difficulty = difficulty
	return &block
}

func (block *Block) GetNowDifficulty(bc *BlockChain, heightNow int, timeStamp uint64) uint64 {
	blockHeight := bc.GetBlockHeight()
	it := bc.NewIterator()
	blockNumberIn5Min := blockHeight
	for {
		block := it.Next()

		// the number of blocks that generated in recent 5 mins
		if timeStamp-block.TimeStamp >= 300 {
			blockNumberIn5Min = heightNow - blockHeight
			break
		}
		if len(block.PrevHash) == 0 {
			fmt.Printf("over")
			blockNumberIn5Min = heightNow - blockHeight
			break
		}
		blockHeight--
	}
	nowDifficulty := block.CalculateDifficulty(blockNumberIn5Min)
	return nowDifficulty
}

func (block *Block) CalculateDifficulty(blockNumberIn5Min int) uint64 {

	return 1234567
}

// 序列化 把一个自定义的数据转化为字节流
// 使用gob包 / binary.Write()
func (block *Block) Serialize() []byte {
	// gob.encode
	var buffer bytes.Buffer

	// 使用gob进行序列化得到字节流
	// 定义一个编码器
	// 使用编码器进行编码
	/*	type Encoder struct {
		mutex      sync.Mutex              // each item must be sent atomically
		w          []io.Writer             // where to send the data
		sent       map[reflect.Type]typeId // which types we've already sent
		countState *encoderState           // stage for writing counts
		freeList   *encoderState           // list of free encoderStates; avoids reallocation
		byteBuf    encBuffer               // buffer for top-level encoderState
		err        error
	}*/
	// func NewEncoder(w io.Writer) *Encoder {
	encoder := gob.NewEncoder(&buffer)
	// func (enc *Encoder) Encode(e interface{}) error {
	err := encoder.Encode(&block)
	if err != nil {
		log.Panic("errors occur when encode")
	}
	//fmt.Printf("value of xiaoMing: %v\n", xiaoMing)
	return buffer.Bytes()
}

func Deserialize(data []byte) Block {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	var block Block
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("errors occur when decode")
	}
	//fmt.Printf("decode result: %v\n",daMing)
	return block
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

// 模拟梅克尔根。只对数据做拼接处理
func (block *Block) MakeMerkelRoot() []byte {
	// 梅克尔根是一个哈希的追加
	// 将交易的哈希值拼接起来 再整体做哈希处理
	/*var info []byte
	for _, tx := range block.Transactions {
		info = append(info, tx.TXID...)
	}
	hash := sha256.Sum256(info)*/
	// func NewMerkleTree(block *Block) []MerkleNode {
	merkleTree := NewMerkleTree(block)
	hash := merkleTree[len(merkleTree)-1].Data
	return hash[:]
}
func NewBlockblock(txs []*Transaction, prevBlockHash []byte) *Block {
	block := Block{
		Version:      00,
		PrevHash:     prevBlockHash,
		MerkelRoot:   []byte{},
		TimeStamp:    uint64(time.Now().Unix()),
		Difficulty:   0, //随便写的无效值
		Nonce:        0, // 无效值
		Hash:         []byte{},
		Transactions: txs,
	}
	block.MerkelRoot = block.MakeMerkelRoot()
	//block.SetHash()
	// 创建一个pow对象

	//difficulty := block.GetNowDifficulty()
	//pow := NewProofOfWork(&block, difficulty)
	// 查找随机数 不停进行哈希运算
	//hash, nonce := pow.Run()
	// 根据挖矿结果 不断对区块数据进行补充
	//block.Hash = hash
	//block.Nonce = nonce
	//block.Difficulty = difficulty
	return &block
}
