package main

import "github.com/gorilla/websocket"

type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

func (c *client) read() {
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			break
		}
		c.room.forward <- msg
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

type room struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	// mapはthread safeじゃないのでchannelを使って操作したい方針
	clients map[*client]bool
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			//welcome
			r.clients[client] = true
		case client := <-r.leave:
			// bye
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// forward messages
			for client := range r.clients {
				select {
				case client.send <- msg:
					// send
				default:
					// failed to send a message
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}
