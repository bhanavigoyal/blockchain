package minerserver

import (
	"encoding/json"
	"fmt"

	"github.com/bhanavigoyal/blockchain/pkg"
)

func NewTransactionHandler(event pkg.Event) error {
	var newTransaction pkg.NewTransactionPayload
	if err := json.Unmarshal(event.Payload, &newTransaction); err != nil {
		return fmt.Errorf("bad payload request: %v", err)
	}

	if err := IsValid(&newTransaction.Transaction); err != nil {
		return fmt.Errorf("invalid Transaction: %v", err)
	}


	//add to mempool
	//add to block and solve hash
	return nil
}

func ReceiveMinedBlockHandler(event pkg.Event) error {
	return nil
}

func SendMinedBlockHandler(event pkg.Event) error {
	return nil
}
