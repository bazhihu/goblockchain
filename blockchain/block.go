package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"goblockchain/merkletree"
	"goblockchain/transaction"
	"goblockchain/utils"
	"time"
)

// 区块
type Block struct {
	Timestamp    int64
	Hash         []byte
	PrevHash     []byte
	Height       int64  // 区块链高度
	Target       []byte // 目标难度值
	Nonce        int64
	Transactions []*transaction.Transaction
	MTree        *merkletree.MerkleTree
}

func (b *Block) SetHash() {
	information := bytes.Join([][]byte{utils.ToHexInt(b.Timestamp), b.PrevHash, b.Target, utils.ToHexInt(b.Nonce), b.BackTrasactionSummary(), b.MTree.RootNode.Data}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}

// 区块创建
func CreateBlock(prevhash []byte, height int64, txs []*transaction.Transaction) *Block {
	block := Block{time.Now().Unix(), []byte{}, prevhash, height, []byte{}, 0, txs, merkletree.CrateMerkleTree(txs)}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return &block
}

// 创始区块
func GenesisBlock(address []byte) *Block {
	//genesisWords := "Hello, blockchain!"
	tx := transaction.BaseTx(address)
	genesis := CreateBlock([]byte("Leo Cao is awesome!"), 0, []*transaction.Transaction{tx})
	genesis.SetHash()
	return genesis
}

// 返回所有交易ID
func (b *Block) BackTrasactionSummary() []byte {
	txIDs := make([][]byte, 0)
	for _, tx := range b.Transactions {
		txIDs = append(txIDs, tx.ID)
	}
	summary := bytes.Join(txIDs, []byte{})
	return summary
}

// -----------------------------------------------------

// 序列化
func (b *Block) Serialize() []byte {
	var (
		res     bytes.Buffer
		encoder = gob.NewEncoder(&res)
		err     = encoder.Encode(b)
	)
	utils.Handle(err)
	return res.Bytes()
}

// 反序列化
func (b *Block) DeSerializeBlock(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	utils.Handle(err)
	return &block
}
