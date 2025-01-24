package centralserver

import (
	"encoding/json"
	"fmt"
	"log"

	pkg "github.com/bhanavigoyal/blockchain/shared"
)

func (m *Manager) NewTransactionHandler(event pkg.Event, client *Client) error {
	var transactionEvent pkg.NewTransactionPayload
	if err := json.Unmarshal(event.Payload, &transactionEvent); err != nil {
		return fmt.Errorf("bad payload request: %v", err)
	}

	data, err := json.Marshal(transactionEvent)

	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent pkg.Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = pkg.EventNewTransaction

	for c := range client.manager.clients {
		c.egress <- outgoingEvent
	}

	return nil
}

func (m *Manager) NewMinedBlockHandler(event pkg.Event, client *Client) error {
	var minedBlockEvent pkg.NewMinedBlockPayload
	if err := json.Unmarshal(event.Payload, &minedBlockEvent); err != nil {
		return fmt.Errorf("bad payload request: %v", err)
	}

	//add vlidation

	// m.Lock()
	// defer m.Unlock()

	if err := m.IsValidBlock(minedBlockEvent.Block); err != nil {
		return fmt.Errorf("error in block validation: %v", err)
	}

	//if validated then broadcast

	data, err := json.Marshal(minedBlockEvent)

	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	outgoingEvent := pkg.Event{
		Payload: data,
		Type:    pkg.EventReceiveNewMinedBlock,
	}

	for c := range client.manager.clients {
		select {
		case c.egress <- outgoingEvent:
		default:
			fmt.Printf("Warning: client %v egress channel full, skipping\n", c)
			//handle skipped clients
		}

	}

	return nil

}

func (m *Manager) SynchronizeMiner(client *Client) {
	m.Lock()
	defer m.Unlock()

	if err := client.connection.WriteJSON(m.blockchain); err != nil {
		log.Printf("Error synchronizing miner: %v", err)
	}
}
