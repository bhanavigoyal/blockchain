package pkg

import (
	"time"
)

type Signature struct {
	R []byte
	S []byte
}

type Transaction struct {
	Message        string
	PublicKey      []byte
	ReceiveAddress []byte
	Amount         string
	Signature      []byte    //change the format -> [R|S]
	Timestamp      time.Time //change the data type
	TxID           []byte
}
