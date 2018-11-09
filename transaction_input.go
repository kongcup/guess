package main

import "bytes"

//TxInput represents a transaction input
type TxInput struct {
	Txid []byte
	Vout int
	Signature []byte
	PubKey []byte
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool  {
	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}