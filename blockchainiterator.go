package main

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockChainInterator struct {
	db *bolt.DB
	// 游标 用于不断索引
	currentHashPointer []byte
}

func (bc *BlockChain) NewIterator() *BlockChainInterator {
	return &BlockChainInterator{
		db: bc.db,
		// 最初指向区块的最后一个区块
		// tail []byte // value of last block's hash
		currentHashPointer: bc.tail,
	}
}

// 迭代器是属于区块链的，Next是迭代器的
// 1. 返回当前区块
// 2. 指针前移
func (it *BlockChainInterator) Next() *Block {
	var block Block
	// func (db *DB) View(fn func(*Tx) error) error {
	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("errors occur when iterator working")
		}
		// 解码
		blockTmp := bucket.Get(it.currentHashPointer)
		block = Deserialize(blockTmp)
		it.currentHashPointer = block.PrevHash
		return nil
	})
	return &block
}
