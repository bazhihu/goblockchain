package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob" // 序列化结构体
	"goblockchain/constcoe"
	"goblockchain/utils"
)

type Transaction struct {
	ID      []byte // 自身hash
	Inputs  []TxInput
	Outputs []TxOutput
}

func (tx *Transaction) TxHash() []byte {
	var (
		encoded bytes.Buffer
		hash    [32]byte
	)
	// 序列化
	var encoder = gob.NewEncoder(&encoded)
	var err = encoder.Encode(tx)
	utils.Handle(err)

	// 加密
	hash = sha256.Sum256(encoded.Bytes())
	return hash[:]
}

func (tx *Transaction) SetID() {
	tx.ID = tx.TxHash()
}

func BaseTx(toaddress []byte) *Transaction {
	var (
		txIn  = TxInput{[]byte{}, -1, []byte{}}
		txOut = TxOutput{constcoe.InitCoin, toaddress}
	)
	tx := Transaction{[]byte("This is the Base Transaction!"), []TxInput{txIn}, []TxOutput{txOut}}
	return &tx
}

func (tx *Transaction) IsBase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].OutIdx == -1
}
