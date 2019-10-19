package main

import "fmt"

func main() {
	// 1. 首先 21万个块 减半
	// 2. 最初 奖励50个
	// 3. 用一个循环 判断 累加
	total := 0.0
	blockInterval := 210000.0 //
	currentReward := 50.0
	for currentReward > 0 {
		//每一个区间内的总量
		amount := blockInterval * currentReward
		currentReward *= 0.5
		// 除法效率低 用等价的乘法
		total += amount
	}
	fmt.Println("Total Number of Bitcoins: 11", total)
}
