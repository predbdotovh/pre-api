package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type triggerAction struct {
	Action string     `json:"action"`
	Row    *sphinxRow `json:"row"`
}

type wsClient struct {
	conn *websocket.Conn
	send chan triggerAction
}

const (
	writeTimeout = 10 * time.Second
	pongWait     = 60 * time.Second
	pingPeriod   = 50 * time.Second
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var backendUpdates chan triggerAction

func websocketUpgrader(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	c := &wsClient{conn: ws, send: make(chan triggerAction)}
	clients[c] = true
	go c.writePump()
	c.readPump()
}

var clients map[*wsClient]bool

func backendPump() {
	backendUpdates = make(chan triggerAction)
	clients = make(map[*wsClient]bool)

	for {
		select {
		case rls := <-backendUpdates:
			for c := range clients {
				select {
				case c.send <- rls:
				default:
					log.Println("CLOSE")
					close(c.send)
					delete(clients, c)
				}
			}
		}
	}
}

func (c *wsClient) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer ticker.Stop()
	for {
		select {
		case rls := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			err := c.conn.WriteJSON(rls)
			if err != nil {
				log.Println("WRITE ERR", err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			err := c.conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println("PING ERR", err)
				return
			}
		}
	}
}

func (c *wsClient) readPump() {
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, _, err := c.conn.NextReader()
		if err != nil {
			log.Println("READ ERR", err)
			break
		}
	}
}
