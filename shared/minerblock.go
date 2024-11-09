package pkg

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

func NewBlockHeader(prevBlockHash []byte) *BlockHeader {
	return &BlockHeader{
		Timestamp:         time.Now().UTC(),
		PreviousBlockHash: prevBlockHash,
		CurrBlockHash:     []byte{},
		Nonce:             0,
		Target:            []byte{0x00, 0x00, 0x00, 0x00},
	}
}

func NewBlockTemplate(blockheader *BlockHeader) *Block {
	return &Block{
		Header:       *blockheader,
		Transactions: []Transaction{},
	}
}

func NewBlock(blockHeader *BlockHeader, transactions Transactions) *Block {
	return &Block{
		Header:       *blockHeader,
		Transactions: transactions,
	}
}

func (b *BlockHeader) ComputeBlockHash() []byte {
	timestamp := []byte(b.Timestamp.String())
	nonce := []byte(strconv.Itoa(b.Nonce))
	headers := bytes.Join([][]byte{timestamp, b.PreviousBlockHash, nonce, b.Target}, []byte{})
	firstHash := sha256.Sum256(headers)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:]
}
