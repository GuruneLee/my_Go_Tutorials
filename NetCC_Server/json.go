package main

import (
	"log"
	"strings"
)

// JSON is the format for communication between client and server
type JSON struct {
	C    *Client `json:"-"`
	Type string  `json:"type"`
	Data Data    `json:"data"`
}

// Data is embedded in JSON
type Data struct {
	ID         string `json:"key"`
	Expression string `json:"expression"`
}

// MakeSendData is sendMe, sendOther를 리턴하는 함수
func (c *Client) MakeSendData(rData JSON, recieveType string) (JSON, JSON) {
	var sendMe JSON
	var sendOther JSON
	switch recieveType {
	case "open":
		log.Printf("Received: %s\n", recieveType)
		sendMe, sendOther = func() (JSON, JSON) {
			IDs := strings.Join(c.server.findIDs(), ",")
			m := JSON{c, "welcome", Data{IDs, "natural"}}
			o := JSON{c, "enter", Data{c.id, "natural"}}
			return m, o
		}()
	case "close":
		log.Printf("Received: %s\n", recieveType)
		sendMe, sendOther = func() (JSON, JSON) {
			m := JSON{c, "bye", Data{c.id, "natural"}}
			o := JSON{c, "exit", Data{c.id, "natural"}}
			return m, o
		}()
	case "exp":
		log.Printf("Received: %s\n", recieveType)
		sendMe, sendOther = func() (JSON, JSON) {
			m := JSON{c, "bye", Data{c.id, "natural"}}
			o := JSON{c, "exit", Data{c.id, "natural"}}
			return m, o
		}()
	}

	return sendMe, sendOther
}
