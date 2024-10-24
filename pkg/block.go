package pkg

import (
	"bytes"
	"crypto/sha256"
	"strconv"
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

func (b* Block) ComputeBlockHash() []byte{
	timestamp := []byte(b.Header.Timestamp.String())
	nonce := []byte(strconv.Itoa(b.Header.Nonce))
	headers := bytes.Join([][]byte{timestamp, b.Header.PreviousBlockHash,nonce }, []byte{})
	firstHash :=sha256.Sum256(headers)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:]
}
