package main

import (
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
func NewBlockChain() *BlockChain {

	/*return &BlockChain{
		blocks: []*Block{genesisBlock},
	}*/
	// 最后一个数据块的哈希
	var lastHash []byte

	// 1. 打开数据库
	// func Open(path string, mode os.FileMode, options *Options) (*DB, error) {
	db, err := bolt.Open(blockChainDB, 0600, nil)
	defer db.Close()
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
			genesisBlock := GenesisBlock()

			// 写数据
			// hash作为key block的字节流作为value
			// func (b *Bucket) Put(key []byte, value []byte) error {
			bucket.Put(genesisBlock.Hash, genesisBlock.toByte())
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
func GenesisBlock() *Block {
	return NewBlock("创世块", []byte{})
}

// 5. 添加区块
func (bc *BlockChain) AddBlock(data string) {
	/*
		// 添加区块数据
		// 更新lastHashKey的value

		// 根据下标获取前区块哈希
		// lastHashKey -> blockN.Hash
		lastBlock := bc.blocks[len(bc.blocks)-1]
		prevHash := lastBlock.Hash
		// 1. 创建新区块
		block := NewBlock(data, prevHash)
		// 2. 添加到区块链数组中
		bc.blocks = append(bc.blocks, block)
	*/

}
