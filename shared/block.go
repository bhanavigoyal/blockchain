package pkg

import (
	"time"
)

type BlockHeader struct {
	Timestamp         time.Time
	PreviousBlockHash []byte
	CurrBlockHash     []byte
	Nonce             int
	Target            []byte
	Height            int
}

type Block struct {
	Header       BlockHeader
	Transactions Transactions
}

type Transactions []Transaction
