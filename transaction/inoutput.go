package transaction

import "bytes"

// 转出的资产值
type TxOutput struct {
	Value     int
	ToAddress []byte // 资产的接收者的地址
}

// 转入的资产值
type TxInput struct {
	TxID        []byte // 前置交易信息
	OutIdx      int
	FromAddress []byte // 资产转出者的地址
}

// 验证收入的地址是否正确
func (in *TxInput) FromAddressRight(address []byte) bool {
	return bytes.Equal(in.FromAddress, address)
}

// 验证支出的地址是否正确
func (out *TxOutput) ToAddressRight(address []byte) bool {
	return bytes.Equal(out.ToAddress, address)
}
