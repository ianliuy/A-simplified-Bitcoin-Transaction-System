package main

// 定义一个wallets结构，保存所有wallet以及它们的地址
type Wallets struct {
	//map[地址][]钱包
	WalletsMap map[string]*Wallet
}

// 创建方法
func NewWallets() *Wallets {
	wallet := NewWallet()
	address := wallet.NewAddress()
	var wallets Wallets
	wallets.WalletsMap = make(map[string]*Wallet)
	wallets.WalletsMap[address] = wallet
	return &wallets

}

// 保存方法，把新建的wallet添加进去

// 读取文件方法，把所有的wallet读取出来
