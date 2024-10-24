package pkg

import (
	"time"
)

type Transaction struct {
	Message         string
	PublicKey       []byte
	ReceieveAddress []byte
	Amount          string
	Signature       []byte
	Timestamp       time.Time
	TxID            []byte
}

type TransactionManager interface{
	ComputeTransactionHash() []byte
}
