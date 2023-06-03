package transaction

import (
	"bytes"
	"encoding/gob"
	"goblockchain/utils"
)

// 转出的资产值
type TxOutput struct {
	Value int
	//ToAddress []byte // 资产的接收者的地址
	HashPubKey []byte // 公钥Hash
}

// 转入的资产值
type TxInput struct {
	TxID   []byte // 前置交易信息
	OutIdx int
	//FromAddress []byte // 资产转出者的地址
	PubKey []byte // 公钥
	Sig    []byte // 签名
}

// 验证收入的地址是否正确
func (in *TxInput) FromAddressRight(address []byte) bool {
	return bytes.Equal(in.PubKey, address)
}

// 验证支出的地址是否正确
func (out *TxOutput) ToAddressRight(address []byte) bool {
	return bytes.Equal(out.HashPubKey, utils.PublicKeyHash(address))
}

// P2PK pay to public key
// 主流区块链系统 用公钥表征Input的地址，用公钥哈希表征Output的地址
// output 使用公钥哈希 能进一步提升区块链系统中交易的匿名性

// 包含output的所有信息
type UTXO struct {
	TxID   []byte
	OutIdx int
	TxOutput
}

// UTXO序列化
func (u *UTXO) Serialize() []byte {
	var res bytes.Buffer
	enCoder := gob.NewEncoder(&res)
	err := enCoder.Encode(u)
	utils.Handle(err)
	return res.Bytes()
}

// UTXO反序列化
func (u *UTXO) DeserializeUTXO(data []byte) *UTXO {
	var utxo UTXO
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&utxo)
	utils.Handle(err)
	return &utxo
}

func DeserializeUTXO(data []byte) *UTXO {
	var utxo UTXO
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&utxo)
	utils.Handle(err)
	return &utxo
}
