package main

import "github.com/gorilla/websocket"

// client express one user doing chat
type client struct {
	// socket is websocket for client
	socket *websocket.Conn
	// send is channel sent message
	send chan []byte
	// room is chatroom the user attend
	room *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
