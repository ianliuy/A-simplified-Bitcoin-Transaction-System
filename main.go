package main

func main() {
	bc := NewBlockChain("19g9h5Xh24Vb92TCwmj1yxL88KwiZ5ytdZ")
	cli := CLI{bc}
	cli.Run()
}
