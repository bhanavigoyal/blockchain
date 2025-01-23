package centralserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	pkg "github.com/bhanavigoyal/blockchain/shared"
	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ErrEventNotSupported = errors.New("this event type is not supported")
)

type EventHandler func(event pkg.Event, client *Client) error

type Manager struct {
	clients  ClientList
	handlers map[string]EventHandler
	sync.RWMutex
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}

	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[pkg.EventNewTransaction] = NewTransactionHandler
	m.handlers[pkg.EventSendNewMinedBlock] = NewMinedBlockHandler
}

func (m *Manager) routeHandler(event pkg.Event, client *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		go func() error {
			if err := handler(event, client); err != nil {
				return err
			}
			return nil
		}()
	} else {
		return ErrEventNotSupported
	}

	return nil
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

// add the handling of disconnection
func (m *Manager) removeClient(client *Client) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		if err := client.connection.Close(); err!=nil{
			fmt.Printf("Error closing connection %v", err)
			return err
		}

		delete(m.clients, client)
		log.Printf("connection closed %v", client.ID)
	}

	return nil

}

func (m *Manager) ServeWs(w http.ResponseWriter, r *http.Request) {
	log.Print("new connection")
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, m)
	m.addClient(client)
	go m.listenToClients(client)
	go client.sendMessage()

}

func (m *Manager) listenToClients(client *Client) {
	defer func() {
		client.manager.removeClient(client)
	}()

	if err := client.connection.SetReadDeadline(time.Now().Add(pkg.PongWait)); err != nil {
		log.Printf("err: %v", err)
		return
	}

	client.connection.SetPongHandler(client.PongHandler)

	//listen to all the incoming events
	for {
		messageType, rawMessage, err := client.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message from %v: %v",client.ID, err)
			} else {
				log.Printf("Error reading Event from %v: %v",client.ID, err)
			}
			if err:=client.manager.removeClient(client); err!=nil{
				log.Printf("error removing client %v: %v",client.ID, err)
			}
			break
		}

		switch messageType {
		case websocket.TextMessage:
			var event pkg.Event
			if err := json.Unmarshal(rawMessage, &event); err != nil {
				log.Printf("Error unmarshaling event from %v: %v",client.ID, err)
				break
			}
			if err = m.routeHandler(event, client); err!= nil{
				log.Printf("error handling new mined block event from %v: %v",client.ID, err)
			}
			

		case websocket.PongMessage:
			var message string
			if err:= json.Unmarshal(rawMessage, &message); err!= nil{
				log.Printf("Error unmarshaling message from %v: %v", client.ID, err)
				continue
			}

			client.connection.SetPongHandler(client.PongHandler)

		}

	}
}
