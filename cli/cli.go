package cli

import (
	"bytes"
	"flag"
	"fmt"
	"goblockchain/blockchain"
	"goblockchain/utils"
	"goblockchain/wallet"
	"log"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct{}

func (cli *CommandLine) printUsage() {
	fmt.Println("Welcome to Leo Cao's tiny blockchain system, usage is as follows:")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("All you need is to first create a wallet.")
	fmt.Println("And then you can use the wallet address to create a blockchain and declare the owner.")
	fmt.Println("Make transactions to expand the blockchain.")
	fmt.Println("In addition, don't forget to run mine function after transatcions are collected.")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("createwallet -refname REFNAME                       ----> Creates and save a wallet. The refname is optional.")
	fmt.Println("walletinfo -refname NAME -address Address           ----> Print the information of a wallet. At least one of the refname and address is required.")
	fmt.Println("walletsupdate                                       ----> Registrate and update all the wallets (especially when you have added an existed .wlt file).")
	fmt.Println("walletslist                                         ----> List all the wallets found (make sure you have run walletsupdate first).")
	fmt.Println("createblockchain -refname NAME -address ADDRESS     ----> Creates a blockchain with the owner you input (address or refname).")
	fmt.Println("balance -refname NAME -address ADDRESS              ----> Back the balance of a wallet using the address (or refname) you input.")
	fmt.Println("blockchaininfo                                      ----> Prints the blocks in the chain.")
	fmt.Println("send -from FROADDRESS -to TOADDRESS -amount AMOUNT  ----> Make a transaction and put it into candidate block.")
	fmt.Println("sendbyrefname -from NAME1 -to NAME2 -amount AMOUNT  ----> Make a transaction and put it into candidate block using refname.")
	fmt.Println("mine                                                ----> Mine and add a block to the chain.")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------")
}

// 创建一个区块链
func (cli *CommandLine) createblockchain(address string) {
	newChain := blockchain.InitBlockChain([]byte(address))
	defer newChain.Database.Close()
	fmt.Println("Finished creating blockchain, and the owner is:", address)
}

// 查看钱包
func (cli *CommandLine) balance(address string) {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	log.Fatalln(chain.Database, address)
	balance, err := chain.FindUTXOs([]byte(address))

	fmt.Printf("Address: %s, Balance:%d , err : %x\n", address, balance, err)
}

// 遍历区块的交易信息
func (cli *CommandLine) getBlockChainInfo() {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()

	// 根据区块链创建迭代器
	iterator := chain.Iterator()
	ogprevhash := chain.BackOgPrevHash()

	for {
		block := iterator.Next()
		fmt.Println("---------------------------------")
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("Transactions: %v\n", block.Transactions)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(block.ValidatePoW()))
		fmt.Println("---------------------------------")
		if bytes.Equal(block.PrevHash, ogprevhash) {
			break
		}
	}
}

// 产生交易信息，并存储到交易信息池中
func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	tx, ok := chain.CreateTransaction([]byte(from), []byte(to), amount)
	if !ok {
		fmt.Println("Failed to create transaction")
		return
	}
	tp := blockchain.CreateTransactionPool()
	tp.AddTransaction(tx)
	tp.SaveFile()
	fmt.Println("Success!")
}

func (cli *CommandLine) mine() {
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()
	chain.RunMine()
	fmt.Println("Finish Mining")
}

// 创建钱包
func (cli *CommandLine) createWallet(refname string) {
	newWallet := wallet.NewWallet()
	newWallet.Save()
	refList := wallet.LoadRefList()
	refList.BindRef(string(newWallet.Address()), refname)
	refList.Save()
	fmt.Println("Success in creating wallet.")
}

// 查看钱包
func (cli *CommandLine) walletInfoRefName(refname string) {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	utils.Handle(err)
	cli.walletInfo(address)
}

// 显示信息
func (cli *CommandLine) walletInfo(address string) {
	wlt := wallet.LoadWallet(address)
	refList := wallet.LoadRefList()
	// 地址公钥
	fmt.Printf("Wallet address: %x\n", wlt.Address())
	fmt.Printf("Public Key:%x\n", wlt.PublicKey)

	// 别名
	fmt.Printf("Reference Name:%s\n", (*refList)[address])
}

// 更新检查
func (cli *CommandLine) walletsUpdate() {
	refList := wallet.LoadRefList()
	refList.Update()
	refList.Save()
	fmt.Println("Success in updating wallets.")
}

// 查看钱包列表
func (cli *CommandLine) walletsList() {
	refList := wallet.LoadRefList()
	for address, _ := range *refList {
		wlt := wallet.LoadWallet(address)
		fmt.Println("-------------------")
		fmt.Printf("Wallet address: %s\n", address)
		fmt.Printf("Public Key :%x\n", wlt.PublicKey)
		fmt.Printf("Reference Name:%s \n", (*refList)[address])
		fmt.Println("-------------------")
		fmt.Println()
	}
}

// 通过别名交易信息
func (cli *CommandLine) sendRefName(fromRefname, toRefname string, amount int) {
	refList := wallet.LoadRefList()
	fromAddress, err := refList.FindRef(fromRefname)
	utils.Handle(err)
	toAddress, err := refList.FindRef(toRefname)
	utils.Handle(err)
	cli.send(fromAddress, toAddress, amount)
}

// 通过别名创建区块链
func (cli *CommandLine) createBlockChainRefName(refname string) {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	utils.Handle(err)
	cli.createblockchain(address)
}

// 通过别名查看余额
func (cli *CommandLine) balanceRefName(refname string) {
	refList := wallet.LoadRefList()
	address, err := refList.FindRef(refname)
	utils.Handle(err)
	cli.balance(address)
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	walletInfoCmd := flag.NewFlagSet("walletinfo", flag.ExitOnError)
	walletsUpdateCmd := flag.NewFlagSet("walletsupdate", flag.ExitOnError)
	walletsListCmd := flag.NewFlagSet("walletslist", flag.ExitOnError)

	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	getBlockChainInfoCmd := flag.NewFlagSet("blockchaininfo", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendByRefNameCmd := flag.NewFlagSet("sendbyrefname", flag.ExitOnError)
	mineCmd := flag.NewFlagSet("mine", flag.ExitOnError)

	createWalletRefName := createWalletCmd.String("refname", "", "The refname of the wallet, and this is optimal") // this line is new
	walletInfoRefName := walletInfoCmd.String("refname", "", "The refname of the wallet")                          // this line is new
	walletInfoAddress := walletInfoCmd.String("address", "", "The address of the wallet")

	createBlockChainOwner := createBlockChainCmd.String("address", "", "The address refer to the owner of blockchain")
	createBlockChainByRefNameOwner := createBlockChainCmd.String("refname", "", "The name refer to the owner of blockchain") // this line is new
	balanceAddress := balanceCmd.String("address", "", "Who needs to get balance amount")
	balanceRefName := balanceCmd.String("refname", "", "Who needs to get balance amount") // this line is new
	sendByRefNameFrom := sendByRefNameCmd.String("from", "", "Source refname")            // this line is new
	sendByRefNameTo := sendByRefNameCmd.String("to", "", "Destination refname")           // this line is new
	sendByRefNameAmount := sendByRefNameCmd.Int("amount", 0, "Amount to send")            // this line is new
	sendFromAddress := sendCmd.String("from", "", "Source address")
	sendToAddress := sendCmd.String("to", "", "Destination address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "walletinfo":
		err := walletInfoCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "walletsupdate":
		err := walletsUpdateCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "walletslist":
		err := walletsListCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "createblockchain":
		err := createBlockChainCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "balance":
		err := balanceCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "blockchaininfo":
		err := getBlockChainInfoCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "mine":
		err := mineCmd.Parse(os.Args[2:])
		utils.Handle(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if createWalletCmd.Parsed() {
		cli.createWallet(*createWalletRefName)
	}

	if walletInfoCmd.Parsed() {
		if *walletInfoAddress == "" {
			if *walletInfoRefName == "" {
				walletInfoCmd.Usage()
				runtime.Goexit()
			} else {
				cli.walletInfoRefName(*walletInfoRefName)
			}
		} else {
			cli.walletInfo(*walletInfoAddress)
		}
	}

	if walletsUpdateCmd.Parsed() {
		cli.walletsUpdate()
	}
	if walletsListCmd.Parsed() {
		cli.walletsList()
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainOwner == "" {
			createBlockChainCmd.Usage()
			runtime.Goexit()
		}
		cli.createblockchain(*createBlockChainOwner)
	}

	if balanceCmd.Parsed() {

		if *balanceAddress == "" {
			if *balanceRefName == "" {
				balanceCmd.Usage()
				runtime.Goexit()
			} else {
				cli.balanceRefName(*balanceRefName)
			}
		} else {
			cli.balance(*balanceAddress)
		}

		cli.balance(*balanceAddress)

	}

	if sendByRefNameCmd.Parsed() {
		if *sendByRefNameFrom == "" || *sendByRefNameTo == "" || *sendByRefNameAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.sendRefName(*sendByRefNameFrom, *sendByRefNameTo, *sendByRefNameAmount)
	}

	if sendCmd.Parsed() {
		if *sendFromAddress == "" || *sendToAddress == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFromAddress, *sendToAddress, *sendAmount)
	}
	if getBlockChainInfoCmd.Parsed() {
		cli.getBlockChainInfo()
	}

	if mineCmd.Parsed() {
		cli.mine()
	}

}
