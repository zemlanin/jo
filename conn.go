package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type connection struct {
	ws   *websocket.Conn
	send chan interface{}
}

func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var message map[string]interface{}
		err := c.ws.ReadJSON(&message)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(message)
		if message["type"] == "NEW_PLAYER" {
			message["playerId"] = "12"
		}
		if message["type"] == "GET_PLAYER" {
			message["player"] = map[string]interface{}{
				"gameId": 2222,
				"name":   "whatever",
				"online": true,
			}
		}
		if message["type"] == "GET_PLAYERS" {
			message["players"] = []map[string]interface{}{
				{
					"gameId": 2222,
					"name":   "whatever",
					"online": true,
				},
				{
					"gameId": 2222,
					"name":   "another",
					"online": false,
				},
			}
		}
		if message["type"] == "GET_GAME_STATE" {
			game_field := map[string]interface{}{
				"x": 0,
				"y": 2,
			}
			message["gameState"] = map[string]interface{}{
				"gameField": game_field,
				"gameId":    2222,
			}
		}
		c.send <- message
	}
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *connection) writeJSON(mt int, payload interface{}) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteJSON(payload)
}

func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.writeJSON(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan interface{}), ws: ws}
	h.register <- c
	go c.writePump()
	c.readPump()
}
