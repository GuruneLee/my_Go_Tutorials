package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type client struct {
	conn  websocket.Conn
	input chan<- map[string]interface{}
}

func (c *client) readInput() {
	//client측에서 주는 메시지(JSON)를 map으로 변환해서 channel로 전달한다
	//channel은 server.newClient에서 server의 channel과 동기화된다
	for {
		mt, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read: ", err)
			break
		}

		var objmap map[string]interface{}
		_ = json.Unmarshal(message, &objmap)

		event := objmap["event"].(string)
		sendData := map[string]interface{}{
			"event":       "res",
			"data":        nil,
			"sender":      c.conn,
			"messageType": mt,
		}

		switch event {
		case "open":
			log.Printf("Received: %s\n", event)
		case "req":
			sendData["data"] = objmap["data"]
			log.Printf("Received: %s\n", event)
		}

		c.input <- sendData //client c의 채널로 senData map전달
	}
}

func (c *client) err(err error, mt int) {
	c.conn.WriteMessage(mt, []byte(err.Error()))
}

func (c *client) msgToClient(msg string, mt int) {
	c.conn.WriteMessage(mt, []byte(msg))
}
