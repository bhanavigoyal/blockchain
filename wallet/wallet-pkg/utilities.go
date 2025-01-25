package wallet

import (
	"crypto/ecdsa"
	"fmt"

	pkg "github.com/bhanavigoyal/blockchain/shared"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func GetAddressFromPrivateKey(generatedPrivateKey ecdsa.PrivateKey) common.Address {

	publicKey := GetPublicKeyFromPrivateKey(generatedPrivateKey)
	pubKey := publicKey.PubKey

	publicKeyEDCSA, ok := pubKey.(*ecdsa.PublicKey)

	if !ok {
		fmt.Printf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	address := crypto.PubkeyToAddress(*publicKeyEDCSA)
	return address
}

func ValidateAddress() {

}

func GetPublicKeyFromPrivateKey(generatedPrivateKey ecdsa.PrivateKey) pkg.PublicKey {
	publicKey := generatedPrivateKey.Public()
	var pubkey pkg.PublicKey
	pubkey.PubKey = publicKey
	return pubkey
}
