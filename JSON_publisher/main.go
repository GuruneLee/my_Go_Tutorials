package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

/*
func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		// func (c *Conn) ReadMessage() (messageType int, p []byte, err error)
		if err != nil {
			log.Println("read:", err)
			break
		}

		var objmap map[string]interface{}
		_ = json.Unmarshal(message, &objmap)
		event := objmap["event"].(string)
		sendData := map[string]interface{}{
			"event": "res",
			"data":  nil,
		}

		switch event {
		case "open":
			log.Printf("Received: %s\n", event)
		case "req":
			sendData["data"] = objmap["data"]
			log.Printf("Received: %s\n", event)
		}

		refineSendData, _ := json.Marshal(sendData)
		err = c.WriteMessage(mt, refineSendData)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
*/
func home(w http.ResponseWriter, r *http.Request) {
	/*
	  type Template
	  : 안전한 HTML조각을 생성하는 text/template의 특별한 형태이다.
	*/

	path := filepath.Join("templates", "index.html")
	// func Join(elem ...string) string
	// : 파라미터들을 각 OS의 구분자를 끼워넣어 연결해준
	tmpl := template.Must(template.ParseFiles(path))
	// func ParseFiles(filenames ...string) (*Template, error)
	// : 새로운 template를 생성하고 filename으로부터 내용을 채워벌임
	// func Must(t *Template, err error) *template
	// : if err가 nil이 아닌 경우,
	// (*Template, error)를 반환하는 함수의 call과 panics를 래핑하는 helper funcion이다
	// -> 변수 초기화에 사용함
	//    ex) var t = template.Must(template.New("name").Parse("html"))
	tmpl.Execute(w, "ws://"+r.Host+"/echo")
	// func (t *Template) Execute(wr io.Writer, data inteface{}) error
	// : 파싱된 template을 특정 data object에 적용하고, output을 wr에다가 쓰는 함수.
}

func main() {

	s := newServer()
	go s.run()

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		serveWs(s, w, r)
	})
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
