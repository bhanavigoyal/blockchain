package centralserver

import (
	"encoding/json"
	"fmt"

	"github.com/bhanavigoyal/blockchain/pkg"
)

func NewTransactionHandler(event pkg.Event, client *Client) error {
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

func NewMinedBlockHandler(event pkg.Event, client *Client) error {
	var minedBlockEvent pkg.NewMinedBlockPayload
	if err := json.Unmarshal(event.Payload, &minedBlockEvent); err != nil {
		return fmt.Errorf("bad payload request: %v", err)
	}

	data, err := json.Marshal(minedBlockEvent)

	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent pkg.Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = pkg.EventReceiveNewMinedBlock

	for c := range client.manager.clients {
		c.egress <- outgoingEvent
	}

	return nil

}
