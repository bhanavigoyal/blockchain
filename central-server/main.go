package main

import (
	"log"
	"net/http"

	centralserver "github.com/bhanavigoyal/blockchain/central-server/central-server-pkg"
)

func main() {
	setupAPI()
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupAPI() {
	manager := centralserver.NewManager()
	http.HandleFunc("/ws", manager.ServeWs)
}