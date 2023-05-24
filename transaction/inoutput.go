package transaction

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
