package main

import (
	"fmt"

	"github.com/bhanavigoyal/blockchain/wallet/wallet-pkg"
)

func main() {
	// address, privateKey := wallet.CreateWallet()
	// fmt.Printf("address: %v \n", address)
	// fmt.Printf("private Key: %v \n", privateKey)

	mnemonic, privateKey, address, err := wallet.GenerateMnemonic()
	if err != nil {
		fmt.Printf("%v:", err)
	}
	fmt.Printf("address: %v \n", address)
	fmt.Printf("mnemonic: %v \n", mnemonic)
	fmt.Printf("privateKey: %v \n", privateKey)
}
