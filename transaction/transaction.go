package transaction

import (
	"bytes"
	"crypto/ecdsa"
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

// 初始化交易信息
func BaseTx(toaddress []byte) *Transaction {
	var (
		txIn  = TxInput{[]byte{}, -1, []byte{}, nil}
		txOut = TxOutput{constcoe.InitCoin, toaddress}
	)
	tx := Transaction{[]byte("This is the Base Transaction!"), []TxInput{txIn}, []TxOutput{txOut}}
	return &tx
}

func (tx *Transaction) IsBase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].OutIdx == -1
}

// 拼装交易新的的交易过程
func (tx *Transaction) PlainCopy() Transaction {
	var (
		inputs  []TxInput
		outputs []TxOutput
	)

	for _, txin := range tx.Inputs {
		inputs = append(inputs, TxInput{txin.TxID, txin.OutIdx, nil, nil})
	}

	for _, txout := range tx.Outputs {
		outputs = append(outputs, TxOutput{txout.Value, txout.HashPubKey})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}
	return txCopy
}

// 对交易信息 填充公钥 input 哈希
func (tx *Transaction) PlainHash(inidx int, prevPubKey []byte) []byte {
	txCopy := tx.PlainCopy()
	txCopy.Inputs[inidx].PubKey = prevPubKey
	return txCopy.TxHash()
}

// 运用公钥 和 私钥签名化
// 分配签名
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey) {
	if tx.IsBase() {
		return
	}

	for idx, input := range tx.Inputs {
		plainhash := tx.PlainHash(idx, input.PubKey)
		signature := utils.Sign(plainhash, privKey)
		tx.Inputs[idx].Sig = signature
	}
}

// 信息交易流程 验证公钥 私钥 和签名
// 分配验证签名
func (tx *Transaction) Verify() bool {
	for idx, input := range tx.Inputs {
		plainhash := tx.PlainHash(idx, input.PubKey)
		if !utils.Verify(plainhash, input.PubKey, input.Sig) {
			return false
		}
	}
	return true
}
