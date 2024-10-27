package pkg

import (
	"bytes"
	"crypto/sha256"
)

func (t *Transaction) ComputeTransactionHash() []byte {
	timestamp := []byte(t.Timestamp.String())

	hashInput := bytes.Join([][]byte{timestamp, t.PublicKey, t.ReceiveAddress, t.Signature}, []byte{})

	firstHash := sha256.Sum256(hashInput)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:]
}

func (t *Transaction) MessageHash() []byte {
	timestamp := []byte(t.Timestamp.String())
	hashInput := bytes.Join([][]byte{timestamp, t.PublicKey, t.ReceiveAddress}, []byte{})

	Z := sha256.Sum256(hashInput)

	return Z[:]
}
