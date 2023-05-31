package constcoe

const (
	Difficulty = 12
	InitCoin   = 1000 //区块链创建时的总币数
	Master     = "Leo Cao"
	Vip1       = "Krad"
	Vip2       = "Exia"

	// 缓冲池，存放一个节点收集到的交易信息
	TransactionPoolFile = "./tmp/transaction_pool.data"

	// 存放区块链数据库的相关地址
	BCPath      = "./tmp/blocks"
	BCFile      = "./tmp/blocks/MANIFEST"
	LastHashKey = "lh"
	BackHashKey = "ogprevhash"

	ChecksumLength = 4
	NetworkVersion = byte(0x00)
	Wallets        = "./tmp/wallets/"
	WalletsRefList = "./tmp/ref_list/"
)
