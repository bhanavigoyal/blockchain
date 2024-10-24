package minerserver

import (
	"bytes"
	"crypto/sha256"
	"strconv"

	"github.com/bhanavigoyal/blockchain/pkg"
)

type BlockMiner struct {
	pkg.Block
}

func (b *BlockMiner) ComputeBlockHash() []byte {
	timestamp := []byte(b.Header.Timestamp.String())
	nonce := []byte(strconv.Itoa(b.Header.Nonce))
	headers := bytes.Join([][]byte{timestamp, b.Header.PreviousBlockHash, nonce}, []byte{})
	firstHash := sha256.Sum256(headers)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:]
}

func (b *BlockMiner) AddTransaction(tx pkg.Transaction) {
	b.Transactions = append(b.Transactions, tx)
}
