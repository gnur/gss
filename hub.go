package main

//import (
//	"strings"
//)

type hub struct {
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
		//case m := <-h.command:
        //    incoming command?
		}
	}
}
