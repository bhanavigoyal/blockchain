package minerserver

import (
	"errors"
	"log"

	"github.com/bhanavigoyal/blockchain/pkg"
	"github.com/gorilla/websocket"
)

var ErrEventNotSupported = errors.New("this event type is not supported")

type EventHandler func(event pkg.Event) error

type Miner struct {
	conn     *websocket.Conn
	handlers map[string]EventHandler
}

func NewMiner(conn *websocket.Conn) *Miner {
	m := &Miner{
		conn:     conn,
		handlers: make(map[string]EventHandler),
	}

	m.setupEventHandlers()
	return m
}

func (m *Miner) setupEventHandlers() {
	m.handlers[pkg.EventNewTransaction] = NewTransactionHandler
	m.handlers[pkg.EventSendNewMinedBlock] = SendMinedBlockHandler
	m.handlers[pkg.EventReceiveNewMinedBlock] = ReceiveMinedBlockHandler
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

		err = m.routeHandler(event)
		if err != nil {
			log.Println("Error handling the event", err)
		}
	}
}
