package centralserver

import (
	"errors"
	"log"
	"net/http"
	"sync"

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
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()

		delete(m.clients, client)
	}
}

func (m *Manager) serveWs(w http.ResponseWriter, r *http.Request) {
	log.Print("new connection")
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, m)
	m.addClient(client)
	go m.listenToClients(client)

}

func (m *Manager) listenToClients(client *Client) {
	defer func() {
		client.manager.removeClient(client)
	}()
	//listen to all the incoming events
	for {
		var event pkg.Event
		err := client.connection.ReadJSON(&event)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break //close connection
		}

		err = m.routeHandler(event, client)
		if err != nil {
			log.Printf("error handling new mined block event: %v", err)
		}
	}
}
