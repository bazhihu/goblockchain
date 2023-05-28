package blockchain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"goblockchain/constcoe"
	"goblockchain/transaction"
	"goblockchain/utils"
	"log"
	"runtime"

	"github.com/dgraph-io/badger"
)

// 区块链
type BlockChain struct {
	// Blocks []*Block
	LastHash []byte
	Database *badger.DB
}

// 构建区块链初始化
// return 包含创始区块
//
//	func CreateBlockChain() *BlockChain {
//		var blockchain BlockChain
//		blockchain.Blocks = append(blockchain.Blocks, GenesisBlock())
//		return &blockchain
//	}
//
// 创建初始化的数据库
func InitBlockChain(address []byte) *BlockChain {
	var lastHash []byte
	if utils.FileExists(constcoe.BCFile) {
		fmt.Println("blockchain already exists")
		runtime.Goexit()
	}
	opts := badger.DefaultOptions(constcoe.BCPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	utils.Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		genesis := GenesisBlock(address)
		fmt.Println("Genesus created")
		// key hash, value 序列化后 serialize
		err = txn.Set(genesis.Hash, genesis.Serialize())
		utils.Handle(err)

		// 设置最近的hash值
		err = txn.Set([]byte(constcoe.LastHashKey), genesis.Hash)
		utils.Handle(err)
		err = txn.Set([]byte(constcoe.BackHashKey), genesis.PrevHash)
		utils.Handle(err)

		lastHash = genesis.Hash
		return err
	})

	utils.Handle(err)
	blockchain := BlockChain{lastHash, db}
	log.Fatalf("%v", blockchain)
	return &blockchain
}

// 读取已有的数据库并加载区块链
func ContinueBlockChain() *BlockChain {
	if utils.FileExists(constcoe.BCFile) == false {
		fmt.Println("No blockchain found, please create on first")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(constcoe.BCPath)
	opts.Logger = nil
	db, err := badger.Open(opts)

	utils.Handle(err)

	err = db.View(func(txn *badger.Txn) error {
		// 先取出最新的hash值
		item, err := txn.Get([]byte(constcoe.LastHashKey))

		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		utils.Handle(err)
		return err
	})

	utils.Handle(err)

	return &BlockChain{lastHash, db}
}

func (bc *BlockChain) AddBlock(newBlock *Block) {
	// newBlock := CreateBlock(bc.Blocks[len(bc.Blocks)-1].Hash, txs)
	// bc.Blocks = append(bc.Blocks, newBlock)
	var lastHash []byte

	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(constcoe.LastHashKey))
		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)
	if !bytes.Equal(newBlock.PrevHash, lastHash) {
		fmt.Println("This block is out of age")
		runtime.Goexit()
	}
	err = bc.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		utils.Handle(err)
		err = txn.Set([]byte(constcoe.LastHashKey), newBlock.Hash)
		bc.LastHash = newBlock.Hash
		return err
	})
	utils.Handle(err)
}

// 区块迭代器
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{chain.LastHash, chain.Database}
}

// 迭代器 Next主函数
func (iterator *BlockChainIterator) Next() *Block {
	var block *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			block = block.DeSerializeBlock(val)
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	iterator.CurrentHash = block.PrevHash
	return block
}

// 迭代器终止器
func (chain *BlockChain) BackOgPrevHash() []byte {
	var ogprevhash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(constcoe.BackHashKey))
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			ogprevhash = val
			return nil
		})

		utils.Handle(err)
		return err
	})

	utils.Handle(err)
	return ogprevhash
}

// 循环区块链中所有区块
// spentTxs 用于记录遍历区块链时那些已经被使用的交易信息的Output
// unSpentTxs就是我们要返回包含指定地址的可用交易信息的切片
func (bc *BlockChain) FindUnspentTransactions(address []byte) []transaction.Transaction {
	var (
		unSpentTxs     []transaction.Transaction
		spentTxs       = make(map[string][]int)
		backOgPrevHash = bc.BackOgPrevHash()
	)

	iter := bc.Iterator()
all:

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		IterOutputs:
			for outIdx, out := range tx.Outputs {

				if spentTxs[txID] != nil {
					for _, spentOut := range spentTxs[txID] {
						if spentOut == outIdx {
							continue IterOutputs
						}
					}
				}

				if out.ToAddressRight(address) {
					unSpentTxs = append(unSpentTxs, *tx)
				}
			}

			if !tx.IsBase() {
				for _, in := range tx.Inputs {
					if in.FromAddressRight(address) {
						fmt.Println("-----------", in.TxID, in.OutIdx)
						inTxID := hex.EncodeToString(in.TxID)
						spentTxs[inTxID] = append(spentTxs[inTxID], in.OutIdx)
					}
				}
			}
		}
		if bytes.Equal(block.PrevHash, backOgPrevHash) {
			break all
		}
	}

	return unSpentTxs
}

func (bc *BlockChain) FindUTXOs(address []byte) (int, map[string]int) {
	unspentOuts := make(map[string]int)

	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {

		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.ToAddressRight(address) {

				accumulated += out.Value
				unspentOuts[txID] = outIdx
				continue Work
			}
		}
	}
	return accumulated, unspentOuts
}

func (bc *BlockChain) FindSpendableOutputs(address []byte, amount int) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {

			if out.ToAddressRight(address) && accumulated < amount {
				fmt.Println(out.Value)
				accumulated += out.Value
				unspentOuts[txID] = outIdx
				if accumulated >= amount {
					break Work
				}
				continue Work
			}
		}
	}

	fmt.Printf("accumulated:%d ; amount: %d \n", accumulated, amount)

	return accumulated, unspentOuts
}

// 打包交易信息
// 产生一条input信息，
// 可能两条outputs 信息，一条是给对方得余额，另一条是我方剩余得余额
func (bc *BlockChain) CreateTransaction(from, to []byte, amount int) (*transaction.Transaction, bool) {
	var (
		inputs  []transaction.TxInput
		outputs []transaction.TxOutput
	)
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		fmt.Println("Not enough coins!")
		return &transaction.Transaction{}, false
	}

	for txid, outidx := range validOutputs {
		txID, err := hex.DecodeString(txid)
		utils.Handle(err)
		input := transaction.TxInput{txID, outidx, from}
		inputs = append(inputs, input)
	}
	outputs = append(outputs, transaction.TxOutput{amount, to})
	if acc > amount {
		fmt.Printf("acc > amount %d, %d\n", acc, amount)
		outputs = append(outputs, transaction.TxOutput{acc - amount, from})
	}

	tx := transaction.Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx, true
}

func (bc *BlockChain) Mine(txs []*transaction.Transaction) {
	// bc.AddBlock(txs)
}
