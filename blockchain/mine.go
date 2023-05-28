package blockchain

import (
	"fmt"
	"goblockchain/utils"
)

// 将收集到的交易信息，存放到区块上，通过Pow共识机制 绑定到区块链上
func (bc *BlockChain) RunMine() {
	transactionPool := CreateTransactionPool()

	candidateBlock := CreateBlock(bc.LastHash, transactionPool.PubTx)
	if candidateBlock.ValidatePoW() {
		bc.AddBlock(candidateBlock)

		err := RemoveTransactionPoolFile()
		utils.Handle(err)
		return
	} else {
		fmt.Println("Block has invalid nonce.")
		return
	}
}
