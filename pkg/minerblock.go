package pkg

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

// Timestamp         time.Time
// 	PreviousBlockHash []byte
// 	CurrBlockHash     []byte
// 	Nonce             int
// 	Difficulty        int

// Header       BlockHeader
// 	Transactions []Transaction

func NewBlockHeader() *BlockHeader {
	return &BlockHeader{
		Timestamp: time.Now().UTC(),
		// PreviousBlockHash: ,
		CurrBlockHash: []byte{},
		Nonce: 0,
	}
}

func NewBlock(blockHeader *BlockHeader, transaction *Transaction) *Block {
	blockHeader.CurrBlockHash = blockHeader.ComputeBlockHash()
	return &Block{
		Header: *blockHeader,
		Transactions: []Transaction{*transaction},
	}
}

func (b *BlockHeader) ComputeBlockHash() []byte {
	timestamp := []byte(b.Timestamp.String())
	nonce := []byte(strconv.Itoa(b.Nonce))
	headers := bytes.Join([][]byte{timestamp, b.PreviousBlockHash, nonce}, []byte{})
	firstHash := sha256.Sum256(headers)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:]
}

func (b *Block) AddTransaction(tx Transaction) {
	b.Transactions = append(b.Transactions, tx)
}
