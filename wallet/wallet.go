package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"fmt"
	"goblockchain/constcoe"
	"goblockchain/utils"
	"io/ioutil"
)

// 钱包功能

// 私钥是一个倍数，公钥是一个点（x,y）。
// 知道基点G（Base Point）和倍数x（私钥）可以在椭圆曲线上计算xG（公钥），而知道G与xG则几乎不可能推测x。
// 使用私钥对信息签名得到的是两个大数，一个是随机数r，一个是计算得到的s。

// 创建椭圆曲线密钥对的生成函数 非对称密钥
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
	err := encoding.Encode(&w)
	fmt.Println(filename, err)
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
