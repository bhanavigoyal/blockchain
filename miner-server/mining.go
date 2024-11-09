package minerserver

import (
	"bytes"
	"time"

	pkg "github.com/bhanavigoyal/blockchain/shared"
)

func (m *Miner) GenerateNewBlock() {
	block := m.blockchain.CreateNewBlock()

	addedTxns := make(map[string]bool)

	startTime := time.Now()

	for len(block.Transactions) < 5 {
		for txId, tx := range m.mempool.transactions {
			if addedTxns[txId] {
				continue
			}

			block.Transactions = append(block.Transactions, *tx)
			addedTxns[txId] = true

			if len(block.Transactions) >= 5 {
				break
			}
		}

		if time.Since(startTime) > time.Minute {
			break
		}

		time.Sleep(10 * time.Second)
	}

	if len(block.Transactions) > 0 {
		newblock := pkg.NewBlock(&block.Header, block.Transactions)
		MineBlock(newblock, m)
	}

}

func MineBlock(b *pkg.Block, m *Miner) {
	for {
		if !bytes.Equal(b.Header.CurrBlockHash[:len(b.Header.Target)], b.Header.Target) {
			b.Header.Nonce += 1
			b.Header.CurrBlockHash = b.Header.ComputeBlockHash()
		} else {
			break
		}
	}

	m.SendMinedBlockHandler(pkg.EventSendNewMinedBlock, *b)

}
