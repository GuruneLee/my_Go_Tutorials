
package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"net/http"
)

func root(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "worng path!\n")
}

func ping(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "pong!\n")
}

func main() {
	conn, err := redis.Dial("tcp", "203.237.53.83:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Do("AUTH", "mypass")
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Do("HMSET", "key", "pings", 0)
	if err != nil {
		log.Fatal(err)
	}

	var i = 0
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		i++
		_, err = conn.Do("HMSET", "key", "pings", i)
		if err != nil {
			log.Fatal(err)
		}
		Pings, err := redis.Int(conn.Do("HGET", "key", "pings"))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(w, "pings: %d", Pings)
	})
	http.HandleFunc("/ping", ping)

	http.ListenAndServe(":8090", nil)
}
