package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"

	pkg "github.com/bhanavigoyal/blockchain/shared"
	"github.com/ethereum/go-ethereum/crypto"
)

func SendTransaction(message string, privateKey ecdsa.PrivateKey, receiveAddress []byte, amount int) error {
	tx := pkg.Transaction{
		Message:        message,
		PublicKey:      crypto.FromECDSAPub(&privateKey.PublicKey),
		ReceiveAddress: receiveAddress,
		Amount:         amount,
		Timestamp:      time.Now(),
	}

	err := SignTransaction(tx, privateKey)

	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	//serialize the txn

	/*

		txJSON, err := json.Marshal(tx)
		if err != nil {
			return fmt.Errorf("failed to serialize transaction: %v", err)
		}
	*/

	//send txn to central server endpoint

	//handle server response

	return nil
}

func SignTransaction(tx pkg.Transaction, privateKey ecdsa.PrivateKey) error {
	txData := fmt.Sprintf("%s%d%x%x", tx.Message, tx.Amount, tx.ReceiveAddress, tx.Timestamp)

	hash := sha256.Sum256([]byte(txData))

	signature, err := ecdsa.SignASN1(rand.Reader, &privateKey, hash[:])
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	tx.Signature = signature
	tx.TxID = tx.ComputeTransactionHash()

	return nil
}

func VerifyTransactionSignature() {

}
