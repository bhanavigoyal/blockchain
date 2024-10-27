package pkg

import (
	"encoding/json"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

const (
	EventNewTransaction       = "new_transaction"
	EventSendNewMinedBlock    = "send_mined_block"
	EventReceiveNewMinedBlock = "receive_mined_block"
)

type NewTransactionPayload struct {
	Transaction Transaction `json:"transaction"`
	From        string      `json:"from"`
}

type NewMinedBlockPayload struct {
	Block Block  `json:"block"`
	From  string `json:"from"`
}
