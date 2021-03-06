package main

import "bytes"

type TxOutput struct {
	Value int
	PubKeyHash []byte
}

// Lock signs the output
func (out *TxOutput) Lock(address []byte)  {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash) - 4]
	out.PubKeyHash = pubKeyHash
}

//IsLockedWithKey checks if the output can be used by the owner of the pubKey
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool  {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

//NewTxOutput create a new TxOutput
func NewTxOutput(value int, address string) *TxOutput  {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}