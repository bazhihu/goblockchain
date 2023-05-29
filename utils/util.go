package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"goblockchain/constcoe"
	"log"
	"os"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)

func ToHexInt(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// 检查文件地址下文件是否存在
func FileExists(fileAddr string) bool {
	if _, err := os.Stat(fileAddr); os.IsNotExist(err) {
		return false
	}
	return true
}

// 公钥转为 公钥哈希
func PublicKeyHash(publicKey []byte) []byte {
	hashedPublicKey := sha256.Sum256(publicKey)
	hasher := ripemd160.New()
	_, err := hasher.Write(hashedPublicKey[:])
	Handle(err)
	publicRipeMd := hasher.Sum(nil)
	return publicRipeMd
}

// 检查位生成函数
func CheckSum(ripeMdHash []byte) []byte {
	firstHash := sha256.Sum256(ripeMdHash)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:constcoe.ChecksumLength]
}

// Base256 和Base58 互相转换
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)
	return []byte(encode)
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	Handle(err)
	return decode
}

// 公钥哈希生成钱包地址
// 将版本号+ 公钥哈希 + 校监sha256值 对称加密 打包成钱包地址
func PubHash2Address(pubKeyHash []byte) []byte {
	networkVersionedHash := append([]byte{constcoe.NetworkVersion}, pubKeyHash...)
	checkSum := CheckSum(networkVersionedHash)

	finalHash := append(networkVersionedHash, checkSum...)
	address := Base58Encode(finalHash)
	return address
}

// 钱包地址 转公钥哈希
// 对称解析 钱包地址，截取中间公钥哈希
func Address2PubHash(address []byte) []byte {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-constcoe.ChecksumLength]
	return pubKeyHash
}
