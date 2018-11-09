package main

import (
	"fmt"
	"log"
)

func (cli *CLI) getBalance(address string)  {
	if !ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc := NewBlockchain(address)
	defer bc.db.Close()

	balance := 0
	pubkeyHash := Base58Decode([]byte(address))
	pubkeyHash = pubkeyHash[1 : len(pubkeyHash) - 4]
	UTXOs := bc.FindUTXO(pubkeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s is %d\n", address, balance)
}
