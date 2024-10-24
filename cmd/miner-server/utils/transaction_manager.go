package minerserver

import (
	"bytes"
	"crypto/sha256"

	"github.com/bhanavigoyal/blockchain/pkg"
)

type Transaction struct {
	pkg.Transaction
}

func (t *Transaction) ComputeTransactionHash() []byte {
	timestamp := []byte(t.Timestamp.String())
	message := []byte(t.Message)
	hashInput := bytes.Join([][]byte{timestamp, t.PublicKey, t.ReceieveAddress, t.Signature, message}, []byte{})

	firstHash := sha256.Sum256(hashInput)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:]
}
