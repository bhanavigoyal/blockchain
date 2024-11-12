package minerserver

import (
	"log"

	"github.com/gorilla/websocket"
)

func main() {
	url := "ws://localhost:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to the server: ", err)
	}
	defer conn.Close()

	mempool := NewMempool()
	miner := NewMiner(conn, mempool)

	go miner.Listen()
	go miner.GenerateNewBlock()
	go miner.sendMessage()
	
	select {}
}
