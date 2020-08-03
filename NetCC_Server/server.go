// server.go
//
// server structure
// func newServer
// func server.run
package main

type server struct {
	clients     map[*Client]bool
	identifiers map[int]*Client
	broadcast   chan map[string]interface{}
	register    chan *Client
	unregister  chan *Client
}

func newServer() *server {
	return &server{
		broadcast:   make(chan map[string]interface{}),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		identifiers: make(map[int]*Client),
	}
}

func (s *server) run() {
	for {
		select {
		// resgister: server에 clients, identifiers 추가하기
		case client := <-s.register:
			s.clients[client] = true
			//식별자 생성
			id := s.makeID()
			//식별자 저장
			s.identifiers[id] = client
			client.id = id

			// resgister: server에서 client, identifiers[client.id] 제거하기
		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				delete(s.identifiers, client.id)
				close(client.send)
			}
		case msg := <-s.broadcast:
			// msg: {conn: , type: , data:{id: , expression: }}
			// conn: &client, nil
			// broadcast에서 읽어온 msg를 선별해서 각자의 client에 뿌림
			for client := range s.clients {
				select {
				case client.send <- msg:
				default: //client.send가 사용불가능 할때, 실행
					close(client.send)
					delete(s.clients, client)
				}
			}
		}
	}
}

func (s *server) makeID() int {
	return len(s.identifiers) + 1
}
