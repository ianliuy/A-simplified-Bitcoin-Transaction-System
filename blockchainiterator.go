package main

import "github.com/boltdb/bolt"

type BlockChainInterator struct {
	db *bolt.DB
	// 游标 用于不断索引
	currentHashPointer []byte
}

func (bc *BlockChain) NewIterator() *BlockChainInterator {
	return &BlockChainInterator{
		db: bc.db,
		// 最初指向区块的最后一个区块
		currentHashPointer: bc.tail,
	}
}
