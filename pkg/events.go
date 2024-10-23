package main

import "encoding/json"

type Event struct{
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

var (
	EventNewTransaction = "new_transaction"
	EventMinedNewBlock = "mined_new_block"
)

type NewTransactionPayload struct{
	Message string `json:"message"`
	From string `json:"from"`
}

type MinedNewBlockPayload struct{
	
}