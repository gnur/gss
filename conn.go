package main

import (
	"code.google.com/p/go.net/websocket"
)

type connection struct {
	ws     *websocket.Conn
	send   chan string
	userid string
}

func (c *connection) reader() {
	var cmd usercmd
	for {
		var message string
		err := websocket.Message.Receive(c.ws, &message)
		if err != nil {
			break
		} else if message == "ka" {

		} else if message[:3] == "id:" {
			c.userid = message[3:]
			h.connect <- c
		} else if message[:3] == "np:" && c.userid != "" {
			cmd.userid = c.userid
			cmd.command = message
			cmd.fromuser = true
			h.command <- cmd
		}
	}
	h.disconnect <- c
}

func (c *connection) writer() {
	for message := range c.send {
		if message == "die" {
			break
		}
		err := websocket.Message.Send(c.ws, message)
		if err != nil {
			break
		}
	}
	h.disconnect <- c
}

func wsHandler(ws *websocket.Conn) {
	c := &connection{send: make(chan string, 256), ws: ws, userid: ""}
	defer func() { h.disconnect <- c }()
	go c.writer()
	c.reader()
}
