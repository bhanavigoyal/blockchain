package minerserver

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	pkg "github.com/bhanavigoyal/blockchain/shared"
	"github.com/gorilla/websocket"
)

var ErrEventNotSupported = errors.New("this event type is not supported")

type EventHandler func(event pkg.Event) error

type Miner struct {
	conn       *websocket.Conn
	handlers   map[string]EventHandler
	mempool    *Mempool
	blockchain *pkg.Blockchain
	egress     chan pkg.Event
	// ingress        chan pkg.Event
	StopMiningChan chan struct{}
}

func NewMiner(conn *websocket.Conn, mempool *Mempool) *Miner {
	m := &Miner{
		conn:           conn,
		handlers:       make(map[string]EventHandler),
		mempool:        mempool,
		StopMiningChan: make(chan struct{}),
		egress:         make(chan pkg.Event, 100),
		//implement blockchain logic for new miner to get current state of blockchain
		blockchain: &pkg.Blockchain{},
	}

	m.synchronizeWithServer(m.conn)
	m.setupEventHandlers()
	return m
}

/*
func (m *Miner) PingHandler(pingMsg string) error {
	log.Printf("Ping received: %v", pingMsg)

	return m.conn.SetReadDeadline(time.Now().Add(pkg.PongWait))
}
*/

func (m *Miner) SendMessage() {

	defer func() {
		m.conn.Close()
		os.Exit(0)
		log.Printf("connection closed")
	}()

	for {
		message, ok := <-m.egress
		if !ok {
			if err := m.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
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
		if err := m.conn.WriteJSON(eventJson); err != nil {
			log.Printf("error sending message: %v", err)
		}
		log.Printf("message sent")
	}
}

func (m *Miner) setupEventHandlers() {
	m.handlers[pkg.EventNewTransaction] = m.NewTransactionHandler
	m.handlers[pkg.EventReceiveNewMinedBlock] = m.ReceiveMinedBlockHandler
}

func (m *Miner) routeHandler(event pkg.Event) error {
	if handler, ok := m.handlers[event.Type]; ok {
		go func() {
			if err := handler(event); err != nil {
				log.Printf("Error handling event %s: %v", event.Type, err)
			}
		}()
	} else {
		return ErrEventNotSupported
	}
	return nil
}

func (m *Miner) Listen() {

	defer func() {
		m.conn.Close()
		os.Exit(0)
	}()

	for {
		messageType, rawMessage, err := m.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			} else {
				log.Printf("Error reading Event: %v", err)
			}
			if err := m.conn.Close(); err != nil {
				log.Printf("Error closing connection: %v", err)
			} else {
				log.Printf("Closed Connection")
				os.Exit(0)
			}
			return
		}

		switch messageType {
		case websocket.TextMessage:
			var event pkg.Event
			if err := json.Unmarshal(rawMessage, &event); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			m.routeHandler(event)

			// case websocket.PingMessage:

			// 	m.conn.SetPingHandler(m.PingHandler)

		}

	}
}

func (m *Miner) synchronizeWithServer(conn *websocket.Conn) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte("SYNC")); err != nil {
		log.Printf("Error requesting synchronization: %v", err)
		return
	}

	_, rawMessage, err := conn.ReadMessage()

	if err != nil {
		log.Printf("Error receiving blockchain from server: %v", err)
		return
	}

	var globalBlockchain pkg.Blockchain
	if err := json.Unmarshal(rawMessage, &globalBlockchain); err != nil {
		log.Printf("Error unmarshaling global blockchain: %v", err)
		return
	}

	m.blockchain = &globalBlockchain
	log.Printf("Miner synchronized with global blockchain")
}
