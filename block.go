package main

import (
	"time"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	Timestamp int64
	Transactions  []*Transaction
	PreBlockHash []byte
	Hash []byte
	Nonce int
}

//NewBlock creates and returns Block
func NewBlock(transactions []*Transaction, preBlockHash []byte) *Block  {
	block := &Block{
		Timestamp:time.Now().Unix(),
		Transactions:transactions,
		PreBlockHash:preBlockHash,
		Hash:[]byte{},
		Nonce:0,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

//NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

//HasTransactions returns a hash of the transactions in the block
func (b *Block) HasTransactions() []byte  {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

//Serialize serializes the block
func (b *Block) Serialize() []byte  {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return  result.Bytes()
}

//DeserializeBlock deserialize a block
func DeserializeBlock(data []byte) *Block  {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}




