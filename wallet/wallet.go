package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"goblockchain/constcoe"
	"goblockchain/utils"
	"io/ioutil"
)

// 钱包功能

// 创建椭圆曲线密钥对的生成函数
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	utils.Handle(err)
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey
}

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// 创建钱包生成函数
func NewWallet() *Wallet {
	privateKey, publicKey := NewKeyPair()
	wallet := Wallet{privateKey, publicKey}
	return &wallet
}

// 先将公钥转成 公钥哈希
// 再将公钥哈希 对称加密生成 钱包地址
func (w *Wallet) Address() []byte {
	pubHash := utils.PublicKeyHash(w.PublicKey)
	return utils.PubHash2Address(pubHash)
}

// 先使用gob时要先注册elliptic.P256()声明elliptic.Curve接口
func (w *Wallet) Save() {
	filename := constcoe.Wallets + string(w.Address()) + ".wlt"
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoding := gob.NewEncoder(&content)
	err := encoding.Encode(w)
	utils.Handle(err)
	err = ioutil.WriteFile(filename, content.Bytes(), 0644)
	utils.Handle(err)
}

// 加载钱包
func LoadWallet(address string) *Wallet {
	filename := constcoe.Wallets + address + ".wlt"
	if !utils.FileExists(filename) {
		utils.Handle(errors.New("no wallet with such address"))
	}
	var w Wallet
	fileContent, err := ioutil.ReadFile(filename)
	utils.Handle(err)
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
	err = decoder.Decode(&w)
	utils.Handle(err)
	return &w
}
