package minerserver

import (
	"bytes"
	"fmt"
	"time"

	pkg "github.com/bhanavigoyal/blockchain/shared"
)

func (m *Miner) GenerateNewBlock() {
	block := m.blockchain.CreateNewBlock()
	fmt.Printf("New Block : %v", block)

	addedTxns := make(map[string]bool)

	startTime := time.Now()

	for len(block.Transactions) < 5 {

		if time.Since(startTime) > time.Minute {
			fmt.Println("Time limit exceeded, stopping txns addition")
			break
		}

		select {
		case <-m.StopMiningChan:
		default:
			fmt.Println("looping txns")
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

			time.Sleep(10 * time.Second)
		}

	}

	if len(block.Transactions) > 0 {
		newblock := pkg.NewBlock(&block.Header, block.Transactions)
		MineBlock(newblock, m)
	} else {
		m.GenerateNewBlock()
	}
}

func MineBlock(b *pkg.Block, m *Miner) {
	for {
		select {
		case <-m.StopMiningChan:
		default:
			if !bytes.Equal(b.Header.CurrBlockHash[:len(b.Header.Target)], b.Header.Target) {
				b.Header.Nonce += 1
				b.Header.CurrBlockHash = b.Header.ComputeBlockHash()
			} else {
				m.SendMinedBlockHandler(pkg.EventSendNewMinedBlock, *b)
				return
			}
		}
	}
}
