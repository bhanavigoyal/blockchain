package minerserver

import (
	"encoding/json"
	"fmt"

	"github.com/bhanavigoyal/blockchain/pkg"
)

func (m *Miner) NewTransactionHandler(event pkg.Event) error {
	var newTransaction pkg.NewTransactionPayload
	if err := json.Unmarshal(event.Payload, &newTransaction); err != nil {
		return fmt.Errorf("bad payload request: %v", err)
	}

	if err := IsValid(&newTransaction.Transaction); err != nil {
		return fmt.Errorf("invalid Transaction: %v", err)
	}

	if err := m.mempool.CheckDoubleSpend(&newTransaction.Transaction); err != nil {
		return fmt.Errorf("double txn: %v", err)
	}

	m.mempool.AddTransaction(&newTransaction.Transaction)

	return nil
}

func (m *Miner) ReceiveMinedBlockHandler(event pkg.Event) error {
	return nil
}

func (m *Miner) SendMinedBlockHandler(event pkg.Event) error {
	return nil
}
