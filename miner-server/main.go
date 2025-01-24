package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	minerserver "github.com/bhanavigoyal/blockchain/miner-server/miner-server-pkg"
	"github.com/gorilla/websocket"
)

func main() {
	url := "ws://localhost:9090/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to the server: ", err)
	}
	defer conn.Close()

	mempool := minerserver.NewMempool()
	miner := minerserver.NewMiner(conn, mempool)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go miner.Listen()
	go miner.GenerateNewBlock()
	go miner.SendMessage()

	<-stopChan

	log.Println("Shutting down miner...")
	close(miner.StopMiningChan)

	log.Println("Shutdown complete.")

}
