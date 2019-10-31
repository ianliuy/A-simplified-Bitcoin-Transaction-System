package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
)

// 演示如何使用ecdsa生成公钥 私钥
// 签名和校验

func main() {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := privateKey.PublicKey
	data := "hello world"
	hash := sha256.Sum256([]byte(data))
	// func Sign(rand io.Reader, priv *PrivateKey, hash []byte) (r, s *big.Int, err error) {
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("publicKey:%v\n", pubKey)
	fmt.Printf("r: %v, len: %v\n", r.Bytes(), len(r.Bytes()))
	fmt.Printf("s: %v, len: %v\n", s.Bytes(), len(s.Bytes()))

	signature := append(r.Bytes(), s.Bytes()...)

	// 1. 定义两个辅助的bigint
	r1 := big.Int{}
	s1 := big.Int{}
	// 2. 拆分signature， 平均分，前半部分给r，后半部分给s
	r1.SetBytes(signature[:len(signature)/2])
	s1.SetBytes(signature[len(signature)/2:])
	// 数据 签名 公钥
	// func Verify(pub *PublicKey, hash []byte, r, s *big.Int) bool {
	res := ecdsa.Verify(&pubKey, hash[:], &r1, &s1)
	fmt.Printf("校验结果：%v\n", res)
}
