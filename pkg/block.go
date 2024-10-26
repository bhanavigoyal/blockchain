package pkg

import (
	"time"
)

type BlockHeader struct {
	Timestamp         time.Time
	PreviousBlockHash []byte
	CurrBlockHash     []byte
	Nonce             int
	Difficulty        int
}

type Block struct {
	Header       BlockHeader
	Transactions []Transaction
}
