package blockchain

// 区块链
type BlockChain struct {
	Blocks []*Block
}

func (bc *BlockChain) AddBlock(data string) {
	newBlock := CreateBlock(bc.Blocks[len(bc.Blocks)-1].Hash, []byte(data))
	bc.Blocks = append(bc.Blocks, newBlock)
}

// 构建区块链初始化
// return 包含创始区块
func CreateBlockChain() *BlockChain {
	var blockchain BlockChain
	blockchain.Blocks = append(blockchain.Blocks, GenesisBlock())
	return &blockchain
}
