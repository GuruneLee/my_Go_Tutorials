package main

type server struct {
	clients   map[*Client]bool
	broadcast chan []byte
	register  chan *Client
}

func newServer() *server {
	return &server{
		broadcast: make(chan []byte),
		register:  make(chan *Client),
		clients:   make(map[*Client]bool),
	}
}

func (s *server) run() {
	for {
		select {
		case client := <-s.register:
			s.clients[client] = true

		case msg := <-s.broadcast:
			for client := range s.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
		}
	}
}
