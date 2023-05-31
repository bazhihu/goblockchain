package blockchain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"goblockchain/transaction"
	"goblockchain/utils"
	"log"
)

// 将收集到的交易信息，存放到区块上，通过Pow共识机制 绑定到区块链上
func (bc *BlockChain) RunMine() {
	transactionPool := CreateTransactionPool()

	if !bc.VerifyTransactions(transactionPool.PubTx) {
		log.Println("falls in transactions verification")
		err := RemoveTransactionPoolFile()
		utils.Handle(err)
		return
	}

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

// 验证交易池中的交易信息 有效性的函数
func (bc *BlockChain) VerifyTransactions(txs []*transaction.Transaction) bool {
	if len(txs) == 0 {
		return true
	}

	spentOutputs := make(map[string]int)
	for _, tx := range txs {
		pubKey := tx.Inputs[0].PubKey
		unspentOuts := bc.FindUnspentTransactions(pubKey)
		inputAmount, OutputAmount := 0, 0

		for _, input := range tx.Inputs {
			if outidx, ok := spentOutputs[hex.EncodeToString(input.TxID)]; ok && outidx == input.OutIdx {
				return false
			}

			ok, amount := isInputRight(unspentOuts, input)
			if !ok {
				return false
			}
			inputAmount += amount
			spentOutputs[hex.EncodeToString(input.TxID)] = input.OutIdx
		}
		for _, output := range tx.Outputs {
			OutputAmount += output.Value
		}

		if inputAmount != OutputAmount {
			return false
		}
		if !tx.Verify() {
			return false
		}
	}
	return true
}

// 验证input 是否正确
func isInputRight(txs []transaction.Transaction, in transaction.TxInput) (bool, int) {
	for _, tx := range txs {
		if bytes.Equal(tx.ID, in.TxID) {
			return true, tx.Outputs[in.OutIdx].Value
		}
	}
	return false, 0
}
