package main

import (
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"strings"
)

type hub struct {
	db         mysql.Conn
	listeners  map[*connection]bool
	command    chan usercmd
	connect    chan *connection
	disconnect chan *connection
}

var h = hub{
	command:    make(chan usercmd),
	connect:    make(chan *connection),
	disconnect: make(chan *connection),
	listeners:  make(map[*connection]bool),
}

func (h *hub) run() {
	h.db = mysql.New("tcp", "", "127.0.0.1:3306", "egdk", "nottheactualpassword", "egdk")
	err := h.db.Connect()
	if err != nil {
		panic(err)
	}
	for {
		select {
		case c := <-h.disconnect:
			if _, ok := h.listeners[c]; ok {
				close(c.send)
				c.ws.Close()
				delete(h.listeners, c)
			}
		case c := <-h.connect:
			found := false
			for val := range h.listeners {
				if val.userid == c.userid {
					val.send <- "die"
					found = true
					break
				}
			}
			if !found || found {
				h.listeners[c] = true
			}
		case m := <-h.command:
			if m.fromuser && m.command[:3] == "np:" {
				commands := strings.Split(m.command[3:], ":")
				//INSERT INTO user_progress (user_id,video_publicid,position) VALUES (1,2,3) ON DUPLICATE KEY UPDATE position=121
				stmt, err := h.db.Prepare("INSERT INTO user_progress (user_hash,video_publicid,position) VALUES (?,?,?) ON DUPLICATE KEY UPDATE position=?")
				if err == nil {
					stmt.Run(m.userid, commands[0], commands[1], commands[1])
				} else {
					h.db = mysql.New("tcp", "", "127.0.0.1:3306", "egdk", "nottheactualpassword", "egdk")
					err := h.db.Connect()
					if err != nil {
						panic(err)
					}
				}
			} else if !m.fromuser {
				for val := range h.listeners {
					if val.userid == m.userid {
						val.send <- m.command
					}
				}
			}
		}
	}
}
