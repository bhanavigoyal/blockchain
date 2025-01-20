package centralserver

import (
	"encoding/json"
	"log"
	"time"

	pkg "github.com/bhanavigoyal/blockchain/shared"
	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn

	manager *Manager

	egress chan pkg.Event
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
	}
}

func (client *Client) sendMessage() {

	ticker := time.NewTicker(pkg.PingInterval)
	defer func() {
		ticker.Stop()
		client.manager.removeClient(client)
	}()

	for {
		select {
		case message, ok := <-client.egress:
			if !ok {
				if err := client.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					// Log that the connection is closed and the reason
					log.Println("connection closed: ", err)
				}
				// Return to close the goroutine
				return
			}
			eventJson, err := json.Marshal(message)
			if err != nil {
				log.Printf("error marshaling event: %v", err)
			}
			if err := client.connection.WriteJSON(eventJson); err != nil {
				log.Printf("error sending message: %v", err)
			}
			log.Printf("message sent")
		case <-ticker.C:
			log.Printf("Ping!")

			if err := client.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("writemsg: %v", err)
				return
			}
		}
	}
}

func (client *Client) PongHandler(pongMsg string) error {
	log.Printf("pong")

	return client.connection.SetReadDeadline(time.Now().Add(pkg.PongWait))
}
