package main

import (
	// "log"
	"jo/players"
)

type hub struct {
	connections map[*connection]bool
	broadcast   chan wsMessage
	register    chan *connection
	unregister  chan *connection
}

var h = hub{
	broadcast:   make(chan wsMessage),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
				players.DisconnectPlayer(c.playerId)
			}
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}
