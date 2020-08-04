package main

import (
	"encoding/json"
	"fmt"
)

type JSON struct {
	Event string `json:"type"`
	Data  data   `json:"data"`
}

type data struct {
	Key        string `json:"key"`
	Expression string `json:"expression"`
}

func main() {
	ex := JSON{
		"exp", data{"idid", "happy"},
	}
	doc, _ := json.Marshal(ex)
	fmt.Println(string(doc))

	var bar JSON
	json.Unmarshal(doc, &bar)
	fmt.Println(bar.Event)

}
