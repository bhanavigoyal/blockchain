package minerserver

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	pkg "github.com/bhanavigoyal/blockchain/shared"
	"github.com/gorilla/websocket"
)

var ErrEventNotSupported = errors.New("this event type is not supported")

type EventHandler func(event pkg.Event) error

type Miner struct {
	conn           *websocket.Conn
	handlers       map[string]EventHandler
	mempool        *Mempool
	blockchain     *pkg.Blockchain
	egress         chan pkg.Event
	ingress        chan pkg.Event
	stopMiningChan chan struct{}
}

func NewMiner(conn *websocket.Conn, mempool *Mempool) *Miner {
	m := &Miner{
		conn:           conn,
		handlers:       make(map[string]EventHandler),
		mempool:        mempool,
		stopMiningChan: make(chan struct{}),
		egress:         make(chan pkg.Event),
		//implement blockchain logic for new miner to get current state of blockchain
	}

	m.setupEventHandlers()
	return m
}

func (m *Miner) PongHandler(pongMsg string) error {
	log.Printf("pong")

	return m.conn.SetReadDeadline(time.Now().Add(pkg.PongWait))
}

func (m *Miner) sendMessage() {

	ticker := time.NewTicker(pkg.PingInterval)
	defer func() {
		ticker.Stop()
		m.conn.Close()
	}()

	for {
		select {
		case message, ok := <-m.egress:
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
		case <-ticker.C:
			log.Printf("Ping!")

			if err := m.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("writemsg: %v", err)
				return
			}

		}
	}
}

func (m *Miner) setupEventHandlers() {
	m.handlers[pkg.EventNewTransaction] = m.NewTransactionHandler
	m.handlers[pkg.EventReceiveNewMinedBlock] = m.ReceiveMinedBlockHandler
}

func (m *Miner) routeHandler(event pkg.Event) error {
	if handler, ok := m.handlers[event.Type]; ok {
		go func() error {
			if err := handler(event); err != nil {
				return err
			}
			return nil
		}()
	} else {
		return ErrEventNotSupported
	}
	return nil
}

func (m *Miner) Listen() {

	defer func() {
		m.conn.Close()
	}()

	if err := m.conn.SetReadDeadline(time.Now().Add(pkg.PongWait)); err != nil {
		log.Printf("err: %v", err)
		return
	}

	m.conn.SetPongHandler(m.PongHandler)

	for {
		var event pkg.Event
		err := m.conn.ReadJSON(&event)
		if err != nil {
			log.Println("Error readiing Event: ", err)
			return
		}
		m.routeHandler(event)

	}
}
