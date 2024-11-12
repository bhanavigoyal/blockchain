package minerserver

import (
	"encoding/json"
	"fmt"

	pkg "github.com/bhanavigoyal/blockchain/shared"
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
	var newBlock pkg.Block

	if err := json.Unmarshal(event.Payload, &newBlock); err != nil {
		return fmt.Errorf("bad payload: %v", err)
	}

	if err := m.IsValidBlock(newBlock); err != nil {
		return fmt.Errorf("error while validating: %v", err)
	}

	return nil
}

func (m *Miner) SendMinedBlockHandler(eventName string, block pkg.Block) error {
	blockJson, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("error marshaling block: %v", err)
	}
	event := &pkg.Event{Type: eventName, Payload: blockJson}
	m.egress <- *event

	return nil
}
