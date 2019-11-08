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
	return NewBlockblock([]*Transaction{coinbase}, []byte{})
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
		block := NewBlock(txs, lastHash, bc)
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
			fmt.Printf("区块数据： %s\n", block.Transactions[0].TXInputs[0].PubKey)
			return nil
		})
		return nil
	})
}

func (bc *BlockChain) GetBlockHeight() int {
	it := bc.NewIterator()
	blockHeight := 0
	for {
		block := it.Next()
		if len(block.PrevHash) == 0 {
			fmt.Printf("over")
			break
		}
		blockHeight++
	}
	return blockHeight
}

// 找到指定地址的所有UTXO
func (bc *BlockChain) FindUTXOs(pubKeyHash []byte) []TXOutput {
	var UTXO []TXOutput

	txs := bc.FindUTXOTransactions(pubKeyHash)
	for _, tx := range txs {
		for _, output := range tx.TXOutputs {
			if bytes.Equal(pubKeyHash, output.PubKeyHash) {
				UTXO = append(UTXO, output)
			}
		}
	}
	return UTXO
}

func (bc *BlockChain) FindNeedUTXOs(senderPubKeyHash []byte, amount float64) (map[string][]uint64, float64) {
	utxos := make(map[string][]uint64)
	var calc float64

	txs := bc.FindUTXOTransactions(senderPubKeyHash)

	for _, tx := range txs {
		for i, output := range tx.TXOutputs {
			//if from == output.PubKeyHash {
			// 变成了两个byte数组的比较
			if bytes.Equal(senderPubKeyHash, output.PubKeyHash) {
				// 3. 比较一下是否满足转账需求
				//   a. 满足的话 直接返回UTXO， calc
				//   b. 不满足的话，继续统计
				if calc < amount {
					// 1. 把UTXO加进来
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], uint64(i))
					// 2. 统计一下当前UTXO总额
					calc += output.Value
					// 加完之后满足条件了
					if calc >= amount {
						fmt.Printf("found satisfied amount: %f\n", calc)
						return utxos, calc
					} else {
						fmt.Printf("当前金额还不满足，当前累计：%f，目标金额：%f\n", calc, amount)
					}
				}
			}
		}
	}
	return utxos, calc
}

func (bc *BlockChain) FindUTXOTransactions(senderPubKeyHash []byte) []*Transaction {
	var txs []*Transaction //存储所有包含utxo的交易
	spentOutputs := make(map[string][]int64)
	it := bc.NewIterator()
	for {
		block := it.Next()
		for _, tx := range block.Transactions {
		OUTPUT:
			for i, output := range tx.TXOutputs {
				if spentOutputs[string(tx.TXID)] != nil {
					for _, j := range spentOutputs[string(tx.TXID)] {
						if int64(i) == j {
							continue OUTPUT
						}
					}
				}
				if bytes.Equal(output.PubKeyHash, senderPubKeyHash) {
					txs = append(txs, tx)
				}
			}
			if !tx.IsCoinbase() {
				for _, input := range tx.TXInputs {
					pubKeyHash := HashPubKey(input.PubKey)
					if bytes.Equal(pubKeyHash, senderPubKeyHash) {
						spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index)
					}
				}
			} else {
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return txs
}

func (bc *BlockChain) FindUTXOTransactionsAAA(senderPubKeyHash []byte) []*Transaction {
	// var UTXO []TXOutput
	var txs []*Transaction //存储所有包含utxo的交易
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
		// 1. iterate block
		block := it.Next()
		// 2. iterate transactions
		for _, tx := range block.Transactions {
			// fmt.Printf("current txid: %x\n", tx.TXID)
			// 3.iterate output(s)
		OUTPUT:
			for i, output := range tx.TXOutputs {
				// fmt.Printf("current idx: %d\n", i)
				// 在这里做一个过滤，将所有削好过的outputs和当前的所将添加output对比一下
				// 如果相同，则跳过。否则继续添加
				// 如果当前的交易ID存在于我们已经表示的map，那么说明这个交易里面有消耗过的output
				if spentOutputs[string(tx.TXID)] != nil {
					for _, j := range spentOutputs[string(tx.TXID)] {
						// []int64{0,1},j:0,1
						if int64(i) == j {
							// 当前准备添加的output已经消耗了，不用再添加了
							continue OUTPUT
						}
					}
				}
				// 如果这个output的地址与目标地址相同，返回utxo数组中
				if bytes.Equal(output.PubKeyHash, senderPubKeyHash) {
					// UTXO = append(UTXO, output)
					//!!!!重点：返回所有包含我的utxo的交易 集合
					txs = append(txs, tx)
					//fmt.Printf("FindUTXO中找到了合适的UTXO，UTXO = append(UTXO, output)\n")
					//fmt.Printf("当前UTXO合并的output.Value为：%v\n", output.Value)
				}
			}
			// 如果当前交易是挖矿交易，不做遍历。跳过
			if !tx.IsCoinbase() {
				// 挖矿交易的id是空，index是-1
				// 遍历input，找到自己花费过的utxo（自己花费过的标示出来）
				for _, input := range tx.TXInputs {
					//判断当前input是否和目标一致
					// 如果相同就加进去
					//if input.PubKey == senderPubKeyHash {
					pubKeyHash := HashPubKey(input.PubKey)
					if bytes.Equal(pubKeyHash, senderPubKeyHash) {
						//spentOutputs := make(map[string][]int64)
						// indexArray := spentOutputs[string(input.TXid)]
						// indexArray = append(indexArray, input.Index)
						spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index)
					}
				}
			} else {
				// fmt.Printf("coinbase, 不做input遍历\n")
			}
		}
		// 3.
		if len(block.PrevHash) == 0 {
			break
			fmt.Printf("区块遍历完成\n")
		}
	}
	return txs
}
