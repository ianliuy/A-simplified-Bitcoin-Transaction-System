package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
)

// 定义一个wallets结构，保存所有wallet以及它们的地址
type Wallets struct {
	//map[地址][]钱包
	WalletsMap map[string]*Wallet
}

// 创建方法
func NewWallets() *Wallets {
	var ws Wallets
	ws.WalletsMap = make(map[string]*Wallet)
	ws.loadFile()
	return &ws
}

func (ws *Wallets) CreatWallet() string {
	wallet := NewWallet()
	address := wallet.NewAddress()
	//var wallets Wallets
	//wallets.WalletsMap = make(map[string]*Wallet)
	ws.WalletsMap[address] = wallet
	ws.saveToFile()
	return address
}

func (ws *Wallets) saveToFile() {
	var buffer bytes.Buffer

	// 因为P256生成的curve类型是一个interface，所以需要跟gob先说一声
	// 在gob注册
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	// func WriteFile(filename string, data []byte, perm os.FileMode) error {
	err = ioutil.WriteFile("wallet.dat", buffer.Bytes(), 0600)
	if err != nil {
		log.Panic(err)
	}
}

// 保存方法，把新建的wallet添加进去
func (ws *Wallets) loadFile() {

	content, err := ioutil.ReadFile("wallet.dat")
	if err != nil {
		log.Panic(err)
	}
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	var wslocal Wallets
	err = decoder.Decode(&wslocal)
	if err != nil {
		log.Panic(err)
	}
	ws.WalletsMap = wslocal.WalletsMap
}

// 读取文件方法，把所有的wallet读取出来
func (ws *Wallets) ListAllAddresses() []string {
	var addresses []string
	for address := range ws.WalletsMap {
		addresses = append(addresses, address)
	}
	return addresses
}
