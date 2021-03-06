package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/miliya12/chat-oreilly/trace"
)

type room struct {
	// forward is channel holding message send to other clients
	forward chan []byte
	// join is channel for clients attempt to join chatroom
	join chan *client
	// leave is channel for clients attempt to leave chatroom
	leave chan *client
	// clients hold all clients join
	clients map[*client]bool
	// tracer gets log of operation executed on chatroom
	tracer trace.Tracer
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// join
			r.clients[client] = true
			r.tracer.Trace("新しいクライアントが参加しました")
		case client := <-r.leave:
			// leave
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("クライアントが退室しました")
		case msg := <-r.forward:
			r.tracer.Trace("メッセージを受信しました: ", string(msg))
			// forward message to all clients
			for client := range r.clients {
				select {
				case client.send <- msg:
					// send message
					r.tracer.Trace(" -- クライアントに送信されました")
				default:
					// failed to send
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(" -- 送信に失敗しました。クライアントをクリーンアップします")
				}
			}
		}
	}
}

const (
	socketBugfferSize = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBugfferSize, WriteBufferSize: messageBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}
