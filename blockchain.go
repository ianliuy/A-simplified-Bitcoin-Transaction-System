package main

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

// 4. 引入区块
type BlockChain struct {
	// 定义一个区块链数组
	// blocks []*Block
	// block Hash -> block.toByte() (转成byte)
	// key是区块的哈希值，value是区块的字节流
	db   *bolt.DB
	tail []byte // value of last block's hash
}

const blockChainDB = "blockChain.db" // database name
const blockBucket = "blockBucket"    // database name

// 5. 定义一个区块链
func NewBlockChain(address string) *BlockChain {

	/*return &BlockChain{
		blocks: []*Block{genesisBlock},
	}*/
	// 最后一个数据块的哈希
	var lastHash []byte

	// 1. 打开数据库
	// func Open(path string, mode os.FileMode, options *Options) (*DB, error) {
	db, err := bolt.Open(blockChainDB, 0600, nil)
	//defer db.Close()
	if err != nil {
		log.Panic("errors occur when opening database")
	}

	// 将要操作数据库
	// func (db *DB) Update(fn func(*Tx) error) error {
	db.Update(func(tx *bolt.Tx) error {
		// 找到抽屉bucket
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			// 没有抽屉 创建
			// func (tx *Tx) CreateBucket(name []byte) (*Bucket, error) {
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("创建bucket（blockBucket）失败")
			}

			// 创建一个创世块并作为第一个区块添加到区块链中
			genesisBlock := GenesisBlock(address)
			fmt.Printf("genesisBlock:%s\n", genesisBlock)
			// hash作为key block的字节流作为value
			// func (b *Bucket) Put(key []byte, value []byte) error {
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			bucket.Put([]byte("LastHashKey"), genesisBlock.Hash)
			lastHash = genesisBlock.Hash
		} else {
			// func (b *Bucket) Get(key []byte) []byte {
			lastHash = bucket.Get([]byte("LastHashKey"))
		}
		return nil
	})
	return &BlockChain{db, lastHash}
}

// 创世块
func GenesisBlock(address string) *Block {
	coinbase := NewCoinbaseTX(address, "genesisBlock")
	// func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// 5. 添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	// 添加区块数据
	// 更新lastHashKey的value
	db := bc.db
	lastHash := bc.tail

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("errors occur: bucket is null")
		}
		block := NewBlock(txs, lastHash)
		// 写数据
		// hash作为key block的字节流作为value
		// func (b *Bucket) Put(key []byte, value []byte) error {
		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte("LastHashKey"), block.Hash)
		lastHash = block.Hash

		// update blockChain in memory, perticularly, its tail
		bc.tail = lastHash

		return nil
	})
	// 根据下标获取前区块哈希
	// 1. 创建新区块
	// 2. 添加到区块链数组中
}

func (bc *BlockChain) PrintChain() {
	blockHeight := 0
	bc.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("blockBucket"))
		// 从第一个Key value进行遍历，到最后一个固定的Key时直接返回
		b.ForEach(func(k, v []byte) error {
			if bytes.Equal(k, []byte("LastHashKey")) {
				return nil
			}
			block := Deserialize(v)
			fmt.Printf("Height:%d\n", blockHeight)
			blockHeight++
			fmt.Printf("版本号： %d\n", block.Version)
			fmt.Printf("前区块哈希： %x\n", block.PrevHash)
			fmt.Printf("梅克尔根： %x\n", block.MerkelRoot)
			fmt.Printf("时间戳： %x\n", block.TimeStamp)
			fmt.Printf("难度值（随便写的）： %d\n", block.Difficulty)
			fmt.Printf("随机数： %d\n", block.Nonce)
			fmt.Printf("区块哈希： %x\n", block.Hash)
			fmt.Printf("区块数据： %s\n", block.Transactions[0].TXInputs[0].Sig)
			return nil
		})
		return nil
	})
}

// 找到指定地址的所有UTXO
func (bc *BlockChain) FindUTXOs(address string) []TXOutput {
	var UTXO []TXOutput
	// 定义一个map保存消费过的utxo，key是output交易过的id，value是交易中索引的数组
	// map[交易id][]int64
	spentOutputs := make(map[string][]int64)

	// 先遍历区块
	// 再遍历交易
	// 再遍历output，找到与自己相关的UTXO（再添加output之前检查是否已经消耗过）
	// 遍历input，找到自己花费过的utxo（自己花费过的标示出来）
	// 创建迭代器
	it := bc.NewIterator()
	for {
		// 1. 遍历区块
		block := it.Next()
		// 2. 遍历交易
		for _, tx := range block.Transactions {
			fmt.Printf("current txid: %x\n", tx.TXID)
			// 3.遍历output
			for i, output := range tx.TXOutputs {
				fmt.Printf("current idx: %d\n", i)
				// 如果这个output的地址与目标地址相同，返回utxo数组中
				if output.PubKeyHash == address {
					UTXO = append(UTXO, output)
				}
			}
			// 遍历input，找到自己花费过的utxo（自己花费过的标示出来）
			for _, input := range tx.TXInputs {
				//判断当前input是否和目标一致
				// 如果相同就加进去
				if input.Sig == address {
					//spentOutputs := make(map[string][]int64)
					indexArray := spentOutputs[string(input.TXid)]
					indexArray = append(indexArray, input.Index)

				}
			}
		}
		// 3.

		if len(block.PrevHash) == 0 {
			break
			fmt.Printf("区块遍历完成")
		}
	}

	return UTXO
}
