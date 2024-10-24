package pkg

import (
	"bytes"
	"crypto/sha256"
	"time"
)

type Transaction struct {
	Message         string
	PublicKey       []byte
	ReceieveAddress []byte
	Amount          string
	Signature       []byte
	Timestamp       time.Time
}

func (t* Transaction) ComputeTransactionHash() []byte{
	timestamp := []byte(t.Timestamp.String())
	message := []byte(t.Message)
	hashInput := bytes.Join([][]byte{timestamp, t.PublicKey, t.ReceieveAddress, t.Signature,message }, []byte{})

	firstHash:= sha256.Sum256(hashInput)
	secondHash :=sha256.Sum256(firstHash[:])

	return secondHash[:]
}
