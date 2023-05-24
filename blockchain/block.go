package blockchain

import (
	"bytes"
	"crypto/sha256"
	"goblockchain/utils"
	"time"
)

// 区块
type Block struct {
	Timestamp int64
	Hash      []byte
	PrevHash  []byte
	Target    []byte // 目标难度值
	Nonce     int64
	Data      []byte
}

func (b *Block) SetHash() {
	information := bytes.Join([][]byte{utils.ToHexInt(b.Timestamp), b.PrevHash, b.Target, utils.ToHexInt(b.Nonce), b.Data}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}

// 区块创建
func CreateBlock(prevhash, data []byte) *Block {
	block := Block{time.Now().Unix(), []byte{}, prevhash, []byte{}, 0, data}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce()
	block.SetHash()
	return &block
}

// 创始区块
func GenesisBlock() *Block {
	genesisWords := "Hello, blockchain!"
	return CreateBlock([]byte{}, []byte(genesisWords))
}
