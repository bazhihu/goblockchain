package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
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
