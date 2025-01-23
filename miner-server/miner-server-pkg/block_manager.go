package minerserver

import (
	"bytes"
	"fmt"

	pkg "github.com/bhanavigoyal/blockchain/shared"
)

func (m *Miner) IsValidBlock(block pkg.Block) error {
	for _, tx := range block.Transactions {
		if err := m.IsValid(&tx); err != nil {
			return fmt.Errorf("invalid txns: %v, %v", tx, err)
		}
	}

	if !bytes.Equal(block.Header.CurrBlockHash[:len(block.Header.Target)], block.Header.Target) {
		return fmt.Errorf("invalid block hash")
	}

	//halt mining of block and
	close(m.StopMiningChan)

	//remove the txn from mempool
	for _, tx := range block.Transactions {
		m.mempool.RemoveTransaction(&tx)
	}

	//start new block creation
	m.StopMiningChan = make(chan struct{})
	go m.GenerateNewBlock()

	//process txns
	m.ProcessingTxns(&block)

	//add to blockchain
	m.blockchain.AddMinedBlock(&block)

	return nil
}
