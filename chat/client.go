package main

import "github.com/gorilla/websocket"
import "time"

// client express one user doing chat
type client struct {
	// socket is websocket for client
	socket *websocket.Conn
	// send is channel sent message
	send chan *message
	// room is chatroom the user attend
	room *room
	// userData holds the information about user
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
