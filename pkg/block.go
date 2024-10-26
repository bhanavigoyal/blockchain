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
}

type Block struct {
	Header       BlockHeader
	Transactions []Transaction
}
