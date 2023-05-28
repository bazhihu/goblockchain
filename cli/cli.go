package cli

import (
	"bytes"
	"flag"
	"fmt"
	"goblockchain/blockchain"
	"goblockchain/utils"
	"log"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct{}

func (cli *CommandLine) printUsage() {
	fmt.Println("Welcome to Leo Cao's tiny blockchain system, usage is as follows:")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("All you need is to first create a blockchain and declare the owner.")
	fmt.Println("And then you can make transactions.")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("createblockchain -address ADDRESS                   ----> Creates a blockchain with the owner you input")
	fmt.Println("balance -address ADDRESS                            ----> Back the balance of the address you input")
	fmt.Println("blockchaininfo                                      ----> Prints the blocks in the chain")
	fmt.Println("send -from FROADDRESS -to TOADDRESS -amount AMOUNT  ----> Make a transaction and put it into candidate block")
	fmt.Println("mine                                                ----> Mine and add a block to the chain")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")

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

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}
func (cli *CommandLine) Run() {

	cli.validateArgs()
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	getBlockChainInfoCmd := flag.NewFlagSet("blockchaininfo", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	mineCmd := flag.NewFlagSet("mine", flag.ExitOnError)

	createBlockChainOwner := createBlockChainCmd.String("address", "", "The address refer to the owner of blockchain")
	balanceAddress := balanceCmd.String("address", "", "Who need to get balance amount")
	sendFromAddress := sendCmd.String("from", "", "Source address")
	sendToAddress := sendCmd.String("to", "", "Destination address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
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

	if createBlockChainCmd.Parsed() {
		if *createBlockChainOwner == "" {
			createBlockChainCmd.Usage()
			runtime.Goexit()
		}
		cli.createblockchain(*createBlockChainOwner)
	}

	if balanceCmd.Parsed() {

		if *balanceAddress == "" {

			balanceCmd.Usage()
			runtime.Goexit()
		}

		cli.balance(*balanceAddress)

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