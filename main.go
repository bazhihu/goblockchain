package main

import (
	"fmt"
	"goblockchain/blockchain"
	"time"
)

func main() {
	// var a = []int{1, 2, 3, 4, 5, 6, 7}
	// b := a[:]
	// fmt.Println(b)
	blockchain := blockchain.CreateBlockChain()
	time.Sleep(time.Second)
	blockchain.AddBlock("After genesis, I have something to say.")
	time.Sleep(time.Second)
	blockchain.AddBlock("Leo Cao is awesome!")
	time.Sleep(time.Second)
	blockchain.AddBlock("I can't wait to follow his github!")
	time.Sleep(time.Second)

	for _, block := range blockchain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("nonce: %d\n", block.Nonce)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Println("Proof of work validation: ", block.ValidatePoW())
	}
}
