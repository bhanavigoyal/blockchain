package minerserver

import (
	"sync"

	"github.com/bhanavigoyal/blockchain/pkg"
)

type Mempool struct {
	transactions map[string]*pkg.Transaction
	sync.Mutex
}

func NewMempool() *Mempool{
	return &Mempool{
		transactions: make(map[string]*pkg.Transaction),
	}
}

func (m *Mempool) AddTransaction(tx *pkg.Transaction) {
	m.Lock()
	defer m.Unlock()

	txId := string(tx.TxID)
	m.transactions[txId] = tx
}

//remove function
