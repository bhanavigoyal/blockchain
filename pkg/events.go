package pkg

import (
	"encoding/json"

	"github.com/bhanavigoyal/blockchain/cmd/central-server"
)
type Event struct{
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, client *centralserver.Client) error

const (
	EventNewTransaction = "new_transaction"
	EventMinedNewBlock = "new_mined_block"
)

type NewTransactionPayload struct{
	Transaction string `json:"transaction"`
	From string `json:"from"`
}

type NewMinedBlockPayload struct{
	Block string `json:"block"`
	From string `json:"from"`
}