package main

import (
	"log"
	"net/http"

	centralserver "github.com/bhanavigoyal/blockchain/central-server/central-server-pkg"
)

func main() {
	setupAPI()
	log.Println("Server started on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func setupAPI() {
	manager := centralserver.NewManager()
	http.HandleFunc("/ws", manager.ServeWs)
}