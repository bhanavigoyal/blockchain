package minerserver

import (
	"time"

	pkg "github.com/bhanavigoyal/blockchain/shared"
)

func (m *Miner) GenerateNewBlock() {
	newblock := m.blockchain.CreateNewBlock()

	addedTxns := make(map[string]bool)

	startTime := time.Now()

	for len(newblock.Transactions) < 5 {
		for txId, tx := range m.mempool.transactions {
			if addedTxns[txId] {
				continue
			}

			newblock.Transactions = append(newblock.Transactions, *tx)
			addedTxns[txId] = true

			if len(newblock.Transactions) >= 5 {
				break
			}
		}

		if time.Since(startTime) > time.Minute {
			break
		}

		time.Sleep(10 * time.Second)
	}

	if len(newblock.Transactions) > 0 {
		pkg.NewBlock(&newblock.Header, newblock.Transactions)
	}

}
