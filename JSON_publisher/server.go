package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type server struct {
	clients []client
	inputs  chan map[string]interface{}
}

func newServer() *server {
	return &server{
		clients: make([]client, 10),
		inputs:  make(chan map[string]interface{}),
	}
}

func (s *server) newClient(conn websocket.Conn) {
	log.Printf("new client has connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn:  conn,
		input: s.inputs,
	}

	s.clients = append(s.clients, *c)

	c.readInput()
}

/*
func (s *server) broadcast() {
	for msg := range s.inputs {
		for c := range s.clients {
			fmt.Println(reflect.TypeOf(c))
		}
	}
}
*/
