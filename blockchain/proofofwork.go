package blockchain

import (
	"bytes"
	"crypto/sha256"
	"goblockchain/constcoe"
	"goblockchain/utils"
	"math"
	"math/big"
)

// Lsh函数就是向左移位，移的越多目标难度值越大，哈希取值落在的空间就更多就越容易找到符合条件的nonce
func (b *Block) GetTarget() []byte {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-constcoe.Difficulty))
	return target.Bytes()
}

func (b *Block) GetBase4Nonce(nonce int64) []byte {
	data := bytes.Join([][]byte{
		utils.ToHexInt(b.Timestamp),
		b.PrevHash,
		utils.ToHexInt(int64(nonce)),
		b.Target,
		b.Data,
	}, []byte{})
	return data
}

func (b *Block) FindNonce() int64 {
	var (
		intHash  big.Int
		inTarget big.Int
		hash     [32]byte
		nonce    int64 = 0
	)
	inTarget.SetBytes(b.Target)

	// nonce 自增 直到由当前nonce得到的区块哈希转化为数值小于目标难度值为止
	for nonce < math.MaxInt64 {
		data := b.GetBase4Nonce(nonce)
		hash = sha256.Sum256(data)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(&inTarget) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce
}

func (b *Block) ValidatePoW() bool {
	var (
		intHash   big.Int
		intTarget big.Int
		hash      [32]byte
	)
	intTarget.SetBytes(b.Target)
	data := b.GetBase4Nonce(b.Nonce)
	hash = sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	if intHash.Cmp(&intTarget) == -1 {
		return true
	}
	return false
}
