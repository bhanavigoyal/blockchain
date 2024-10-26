package minerserver

import (
	"fmt"

	"github.com/bhanavigoyal/blockchain/pkg"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func IsValid(t *pkg.Transaction) error {
	z := t.MessageHash()
	if !secp256k1.VerifySignature(t.PublicKey, z, t.Signature) {
		return fmt.Errorf("failed to verify secp256k1 signature")
	}

	txid := t.ComputeTransactionHash()
	t.TxID = txid

	return nil
}
