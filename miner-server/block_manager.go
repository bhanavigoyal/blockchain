package minerserver

import (
	"bytes"
	"fmt"

	pkg "github.com/bhanavigoyal/blockchain/shared"
)

func IsValidBlock(block pkg.Block) error {
	for _, tx := range block.Transactions {
		if err := IsValid(&tx); err != nil {
			return fmt.Errorf("invalid txns: %v, %v", tx, err)
		}
	}

	if !bytes.Equal(block.Header.CurrBlockHash[:len(block.Header.Target)], block.Header.Target) {
		return fmt.Errorf("invalid block hash")
	}
	//remove the txn from mempool
	//halt mining of block and start new block creation
	//add to blockchain

	return nil
}
