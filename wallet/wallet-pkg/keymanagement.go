package wallet

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

func CreateWallet() (common.Address, *ecdsa.PrivateKey) {
	generatedPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Printf("error generating private key: %v", err)
	}

	// privateKeyBytes := crypto.FromECDSA(generatedPrivateKey)
	// privateKey := hexutil.Encode(privateKeyBytes[2:])

	address := GetAddressFromPrivateKey(*generatedPrivateKey)

	return address, generatedPrivateKey

}

func GenerateMnemonic() (string, *ecdsa.PrivateKey, common.Address, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", &ecdsa.PrivateKey{}, common.Address{}, fmt.Errorf("error generating entropy: %v", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", &ecdsa.PrivateKey{}, common.Address{}, fmt.Errorf("error generating mnemonic: %v", err)
	}

	privateKey, addresss, err := CreateWalletFromMnemonic(mnemonic)

	if err != nil {
		return mnemonic, privateKey, addresss, fmt.Errorf("%v", err)
	}

	return mnemonic, privateKey, addresss, nil

}

func CreateWalletFromMnemonic(mnemonic string) (*ecdsa.PrivateKey, common.Address, error) {
	seed := bip39.NewSeed(mnemonic, "")

	generatedPrivateKey, err := crypto.ToECDSA(seed[:32])
	if err != nil {
		return &ecdsa.PrivateKey{}, common.Address{}, fmt.Errorf("error generating private key from seed: %v", err)
	}

	// privateKeyBytes := crypto.FromECDSA(generatedPrivateKey)
	// privateKey := hexutil.Encode(privateKeyBytes[2:])

	address := GetAddressFromPrivateKey(*generatedPrivateKey)

	return generatedPrivateKey, address, nil

}

func ImportWallet() {

}

func ExportPrivateKey() {

}
