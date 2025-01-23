package centralserver

import (
	"encoding/json"
	"log"
	"time"

	pkg "github.com/bhanavigoyal/blockchain/shared"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	ID         string
	connection *websocket.Conn
	manager    *Manager
	egress     chan pkg.Event
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		ID: uuid.NewString(),
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
			//ok -> false if egress channel is closed
			if !ok {
				if err := client.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					// Log that the connection is closed and the reason
					log.Printf("connection closed for %v: %v ",client.ID, err)
				}
				// Return to close the goroutine
				return
			}
			eventJson, err := json.Marshal(message)
			if err != nil {
				log.Printf("error marshaling event from %v: %v",client.ID, err)
				return
			}
			if err := client.connection.WriteJSON(eventJson); err != nil {
				log.Printf("error sending message to %v: %v",client.ID, err)
				return
			}
			log.Printf("message sent to : %v", client.ID)
		case <-ticker.C:
			log.Printf("Sending Ping! to %v", client.ID)

			if err := client.connection.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(pkg.PongWait)); err != nil {
				log.Printf("Error sending Ping to %v: %v",client.ID, err)
				return
			}
		}
	}
}

func (client *Client) PongHandler(pongMsg string) error {
	log.Printf("Pong received from %v: %v",client.ID, pongMsg)

	return client.connection.SetReadDeadline(time.Now().Add(pkg.PongWait))
}
