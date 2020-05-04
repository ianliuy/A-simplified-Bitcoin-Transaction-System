package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

type Wallet struct {
	Private *ecdsa.PrivateKey
	PubKey  []byte
}

// creat New Wallet with privateKey and publicKey
func NewWallet() *Wallet {
	// 生成曲线
	curve := elliptic.P256()
	// 创建私钥
	privateKey, _ := ecdsa.GenerateKey(curve, rand.Reader)
	//if err != nil {
	//	fmt.Printf("errors occur when Generate Key\n")
	//	log.Panic(err)
	//}
	pubKeyOrig := privateKey.PublicKey
	pubKey := append(pubKeyOrig.X.Bytes(), pubKeyOrig.Y.Bytes()...)
	/*	fmt.Printf("PriviteKeyNum: %v\n", privateKey.D)
		fmt.Printf("PublicKeyNum1: %v\n", privateKey.PublicKey.X)
		fmt.Printf("PublicKeyNum2: %v\n", privateKey.PublicKey.Y)*/
	return &Wallet{Private: privateKey, PubKey: pubKey}
}

// 生成地址

func (w *Wallet) NewAddress() string {
	pubKey := w.PubKey
	rip160HashValue := HashPubKey(pubKey)
	version := byte(00)
	// 拼接version
	payload := append([]byte{version}, rip160HashValue...)
	// checksum
	checkCode := CheckSum(payload)
	// 25字节数据
	payload = append(payload, checkCode...)

	// go语言库：btcd 这个是go语言实现的比特币全节点源码
	address := base58.Encode(payload)
	return address
}

func HashPubKey(data []byte) []byte {
	hash := sha256.Sum256(data)

	// 理解为编码器
	rip160hasher := ripemd160.New()
	_, err := rip160hasher.Write(hash[:])
	if err != nil {
		fmt.Println("bbbbb")
		log.Panic(err)
	}
	// 返回rip160的哈希结果
	rip160HashValue := rip160hasher.Sum(nil)
	return rip160HashValue
}
func CheckSum(data []byte) []byte {
	// 两次sha256
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])
	// 前4字节 校验码
	checkCode := hash2[:4]
	return checkCode
}
