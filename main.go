package main

import (
	"goblockchain/blockchain/cli"
	"os"
)

func main() {

	defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Run()

	// for _, block := range blockchain.Blocks {
	// 	fmt.Printf("Timestamp: %d\n", block.Timestamp)
	// 	fmt.Printf("Hash: %x\n", block.Hash)
	// 	fmt.Printf("PrevHash: %x\n", block.PrevHash)
	// 	fmt.Printf("nonce: %d\n", block.Nonce)
	// 	fmt.Printf("Data: %s\n", block.Data)
	// 	fmt.Println("Proof of work validation: ", block.ValidatePoW())
	// }

	// pl := func(chain *blockchain.BlockChain, i int) {
	// 	property, _ := chain.FindUTXOs([]byte(constcoe.Master))
	// 	fmt.Println("Balance of master:", property)

	// 	property, _ = chain.FindUTXOs([]byte(constcoe.Vip1))
	// 	fmt.Println("Balance of VIP1:", property)

	// 	property, _ = chain.FindUTXOs([]byte(constcoe.Vip2))
	// 	fmt.Println("Balance of VIP2:", property)
	// 	fmt.Printf("-----------%d----------\n", i)
	// }

	// txPool := make([]*transaction.Transaction, 0)
	// var (
	// 	tempTx   *transaction.Transaction
	// 	ok       bool
	// 	property int
	// 	chain    = blockchain.CreateBlockChain()
	// )
	// property, _ = chain.FindUTXOs([]byte(constcoe.Master))
	// fmt.Println("Balance of Leo Cao:", property)

	// tempTx, ok = chain.CreateTransaction([]byte(constcoe.Master), []byte(constcoe.Vip1), 100)
	// if ok {
	// 	txPool = append(txPool, tempTx)
	// }

	// // fmt.Println(txPool)
	// chain.Mine(txPool)
	// txPool = make([]*transaction.Transaction, 0)
	// pl(chain, 1)
	// // fmt.Println(txPool)

	// //二次交易
	// tempTx, ok = chain.CreateTransaction([]byte(constcoe.Vip1), []byte(constcoe.Vip2), 200)
	// if ok {
	// 	txPool = append(txPool, tempTx)
	// }
	// pl(chain, 2)
	// fmt.Println(txPool)

	// // 三次交易
	// tempTx, ok = chain.CreateTransaction([]byte(constcoe.Vip1), []byte(constcoe.Vip2), 50)
	// if ok {
	// 	txPool = append(txPool, tempTx)
	// 	fmt.Println(tempTx.Outputs)
	// }

	// pl(chain, 3)
	// fmt.Println(txPool)

	// // 四次交易
	// tempTx, ok = chain.CreateTransaction([]byte(constcoe.Master), []byte(constcoe.Vip2), 100)
	// if ok {
	// 	txPool = append(txPool, tempTx)
	// }
	// pl(chain, 4)

	// chain.Mine(txPool)
	// txPool = make([]*transaction.Transaction, 0)

	// pl(chain, 5)

	// // 五次交易
	// tempTx, ok = chain.CreateTransaction([]byte(constcoe.Master), []byte(constcoe.Vip1), 50)
	// if ok {
	// 	txPool = append(txPool, tempTx)
	// 	fmt.Println(tempTx.Outputs)
	// }

	// // 六次交易
	// tempTx, ok = chain.CreateTransaction([]byte(constcoe.Vip1), []byte(constcoe.Vip2), 100)
	// if ok {
	// 	txPool = append(txPool, tempTx)
	// 	fmt.Println(tempTx.Outputs)
	// }
	// chain.Mine(txPool)
	// pl(chain, 6)

	// property, _ = chain.FindUTXOs([]byte(constcoe.Master))
	// fmt.Println("Balance of master:", property)

	// for _, block := range chain.Blocks {
	// 	fmt.Printf("Timestamp: %d\n", block.Timestamp)
	// 	fmt.Printf("Hash: %x\n", block.Hash)
	// 	fmt.Printf("PrevHash: %x\n", block.PrevHash)
	// 	fmt.Printf("nonce: %d\n", block.Nonce)
	// 	fmt.Println("Proof of work validation: ", block.ValidatePoW())
	// }

	// // bug
	// tempTx, ok = chain.CreateTransaction([]byte(constcoe.Vip1), []byte(constcoe.Vip2), 30)
	// if ok {
	// 	txPool = append(txPool, tempTx)
	// }
	// tempTx, ok = chain.CreateTransaction([]byte(constcoe.Vip1), []byte(constcoe.Master), 30)
	// if ok {
	// 	txPool = append(txPool, tempTx)
	// }

	// chain.Mine(txPool)
	// txPool = make([]*transaction.Transaction, 0)

	// for _, block := range chain.Blocks {
	// 	fmt.Printf("Timestamp: %d\n", block.Timestamp)
	// 	fmt.Printf("Hash: %x\n", block.Hash)
	// 	fmt.Printf("PrevHash: %x\n", block.PrevHash)
	// 	fmt.Printf("nonce: %d\n", block.Nonce)
	// 	fmt.Println("Proof of work validation: ", block.ValidatePoW())
	// }

	// property, _ = chain.FindUTXOs([]byte(constcoe.Master))
	// fmt.Println("Balance of master:", property)
	// property, _ = chain.FindUTXOs([]byte(constcoe.Vip1))
	// fmt.Println("Balance of VIP1:", property)
	// property, _ = chain.FindUTXOs([]byte(constcoe.Vip2))
	// fmt.Println("Balance of VIP2:", property)
}
