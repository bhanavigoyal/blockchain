package minerserver

import (
	"fmt"

	pkg "github.com/bhanavigoyal/blockchain/shared"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func (m *Miner) IsValid(t *pkg.Transaction) error {
	//check for the balance of the wallet
	if m.blockchain.Balances[string(t.PublicKey)] < t.Amount {
		return fmt.Errorf("not enough balance")
	}

	//verify signature
	z := t.MessageHash()
	if !secp256k1.VerifySignature(t.PublicKey, z, t.Signature) {
		return fmt.Errorf("failed to verify secp256k1 signature")
	}

	//compute hash
	txid := t.ComputeTransactionHash()
	t.TxID = txid

	return nil
}

func (m *Miner) ProcessingTxns(block *pkg.Block) error {
	for _, tx := range block.Transactions {
		if err := m.IsValid(&tx); err == nil {
			m.blockchain.Balances[string(tx.PublicKey)] -= tx.Amount
			m.blockchain.Balances[string(tx.ReceiveAddress)] += tx.Amount
		}
	}

	return nil
}
