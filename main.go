package main

import (
	"code.google.com/p/go.net/websocket"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type usercmd struct {
	userid   string
	command  string
	fromuser bool
}

func apiHandler(c http.ResponseWriter, req *http.Request) {
	var err error
	var cmd usercmd
    if match, er := regexp.MatchString("egdk|0x8|erwin", req.Referer()); er == nil && match {
        c.Header().Set("access-control-allow-origin", "*")
    }
	cmd.fromuser = false
	cmd.userid, cmd.command, err = parseUrl(req.URL.Path)
	fmt.Println(cmd.userid)
	if err == nil {
		h.command <- cmd
	} else {
		fmt.Print(err)
		fmt.Print(cmd.userid)
	}
}

func main() {
	go h.run()
	http.HandleFunc("/", apiHandler)
	http.Handle("/ws", websocket.Handler(wsHandler))
	err := http.ListenAndServe(":9010", nil)
	if err != nil {
		fmt.Println("het is niet gelukt.. helaas")
	}
}

func parseUrl(s string) (userid string, function string, err error) {
	userid, function, err = "", "", nil
	parts := strings.Split(s[1:], "/")
	if len(parts) <= 1 {
		err = errors.New("incorrect path")
		return
	}
	if match, er := regexp.MatchString("^[a-f0-9]{40}$", parts[0]); er == nil && match {
		userid = parts[0]
	} else {
		err = errors.New("incorrect userid: '" + parts[0] + "'")
		return
	}
	function = parts[1]
	if len(parts) == 2 || parts[2] == "" {
		function = function + "()"
		return
	}
	paramstring := ""
	for _, param := range parts[2:] {
		if match, er := regexp.MatchString("/^[0-9]*$/", param); er == nil && match {
			paramstring = paramstring + ", " + param
		} else if len(param) > 0 {
			paramstring = paramstring + ", '" + param + "'"
		}
	}
	function = function + "(" + paramstring[2:] + ")"
	return
}
