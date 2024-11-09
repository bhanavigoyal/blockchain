package minerserver

import (
	"errors"
	"log"

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
}

func NewMiner(conn *websocket.Conn, mempool *Mempool) *Miner {
	m := &Miner{
		conn:     conn,
		handlers: make(map[string]EventHandler),
		mempool:  mempool,
		//implement blockchain logic for new miner to get current state of blockchain
	}

	m.setupEventHandlers()
	return m
}

func (m *Miner) setupEventHandlers() {
	m.handlers[pkg.EventNewTransaction] = m.NewTransactionHandler
	m.handlers[pkg.EventReceiveNewMinedBlock] = m.ReceiveMinedBlockHandler
}

func (m *Miner) routeHandler(event pkg.Event) error {
	if handler, ok := m.handlers[event.Type]; ok {
		return handler(event)
	}

	return ErrEventNotSupported
}

func (m *Miner) Listen() {
	for {
		var event pkg.Event
		err := m.conn.ReadJSON(&event)
		if err != nil {
			log.Println("Error readiing Event: ", err)
			return
		}

		//concurrent event handling??
		go func(event pkg.Event) {
			err = m.routeHandler(event)
			if err != nil {
				log.Println("Error handling the event", err)
			}
		}(event)
	}
}
