// client.go
//
// Client structure
// func (client) readPump
// func serveWs
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is...
type Client struct {
	server *server
	conn   *websocket.Conn
	id     int
	send   chan map[string]interface{}
}

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		// message 받아오기
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// 받아온 Json파일을 parsing
		//
		// {type: , data: {id: , expression: }}
		// type은 open, close, exp
		var objmap map[string]interface{}
		_ = json.Unmarshal(message, &objmap)
		recieveType := objmap["type"].(string)

		// 보낼 데이터 정의하기
		//
		// sendMe -> 요청한놈
		// sendOther -> 다른놈
		// {conn:, type: , data: {id: , expression: }}
		sendMe := map[string]interface{}{
			// type: welcome, bye, exp
			"conn": c,
			"type": nil,
			"data": objmap["data"],
		}
		sendOther := map[string]interface{}{
			// type: enter, exit, exp
			"conn": c,
			"type": nil,
			"data": objmap["data"],
		}

		// 받은 type따라 sendData 채우기
		switch recieveType {
		case "open":
			log.Printf("Received: %s\n", recieveType)
      sendMe["type"] = "welcome"
      sendOther["type"] = "Enter"
      sendMe["data"]["id"] =
      sendOther["data"]["id"] =
		case "close":
			log.Printf("Received: %s\n", recieveType)
      sendMe["type"] = "bye"
      sendOther["type"] = "Exit"
      sendMe["data"] =
      sendOther["data"] =
		case "exp":
			log.Printf("Received: %s\n", recieveType)
      sendMe["type"] = "exp"
      sendOther["type"] = "exp"
      sendMe["data"] = objmap["data"]
      sendOther["data"] = objmap["data"]
		}

		c.server.broadcast <- sendMe
		c.server.broadcast <- sendOther
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			// routing
			var refineSendData []byte
			if c.rightMSG(message) {
				//parsing
				msg := map[string]interface{}{
					"type": message["type"],
					"data": message["data"],
				}
				//Marshal
				refineSendData, _ = json.Marshal(msg)
			}

			w.Write(refineSendData)

			/*// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			*/
			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) rightMSG(msg map[string]interface{}) bool {
	ans := false
	conn := msg["conn"]
	t := msg["type"]
	bar1 := (conn == c) && ((t == "welcome") || (t == "bye"))
	bar2 := (conn != c) && ((t == "enter") || (t == "exit"))

	if (t == "exp") || bar1 || bar2 {
		ans = true
	}

	return ans
}

func serveWs(hub *server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{server: hub, conn: conn, id: -1, send: make(chan map[string]interface{})}
	client.server.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
