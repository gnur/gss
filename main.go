package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net/http"
)

type usercmd struct {
	userid   string
	command  string
	fromuser bool
}

func main() {
	go h.run()
	http.Handle("/", websocket.Handler(wsHandler))
	err := http.ListenAndServe(":9010", nil)
	if err != nil {
		fmt.Println("het is niet gelukt.. helaas")
	}
}
