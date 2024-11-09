package minerserver

import (
	"fmt"
	"sync"

	pkg "github.com/bhanavigoyal/blockchain/shared"
)

type Mempool struct {
	transactions map[string]*pkg.Transaction
	sync.Mutex
}

func NewMempool() *Mempool {
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

func (m *Mempool) CheckDoubleSpend(tx *pkg.Transaction) error {
	if _, ok := m.transactions[string(tx.TxID)]; ok {
		return fmt.Errorf("transaction with %v already present", tx.TxID)
	} else {
		return nil
	}
}

//remove function

