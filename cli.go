package main

import (
	"fmt"
	"flag"
	"os"
	"log"
)

//CLI responsible for processing command line arguments
type CLI struct {}

func (cli *CLI) printUsage()  {
	fmt.Println("Usage:")
	fmt.Println("   createblockchain -address ADDRESS - Create a block chain and send genesis reward to the ADDRESS")
	fmt.Println("   createwallet - Generate a new key-pair and save it into the wallet file")
	fmt.Println("   getbalance -address ADDRESS - Get Balance of ADDRESS")
	fmt.Println("   listaddresses - Lists all addresses from the wallet file")
	fmt.Println("   printchain - print all blocks of the blockchain")
	fmt.Println("   send -from FROM -to TO -amount Amount - Send AMOUNT of coins  from FROM address to TO")
	fmt.Println("   guess -pwd - guess address ,then send mail")
}

func (cli *CLI) validateArgs()  {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

//Run parses command line arguments and processes commands
func (cli *CLI) Run()  {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd  := flag.NewFlagSet("printchain", flag.ExitOnError)
	guessCmd  := flag.NewFlagSet("guess", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to getBalance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	pwd := guessCmd.String("pwd", "", "email password")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(nil)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "guess":
		err := guessCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses()
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
	if guessCmd.Parsed() {
		cli.guess(*pwd)
	}
}

