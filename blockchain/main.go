package main

import (
	"blockchain/btc"
)

func main() {
	bc := btc.NewBlockchain()
	defer bc.Close()

	cli := btc.CLI{Bc: bc}
	cli.Run()
}
