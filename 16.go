package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"golang.org/x/net/websocket"
)

const LAUNCH_FAILED = 1
const FILE_READ = 2
const BAD_TEMPLATE = 3

const ADDRESS = ":3000"

type PageConfiguration struct {
	URL      string
	Commands []string
}

func main() {
	p := PageConfiguration{"/socket", []string{"A", "B", "C"}}

	html, e := template.ParseFiles(BaseName() + ".html")
	Abort(FILE_READ, e)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		Abort(BAD_TEMPLATE, html.Execute(w, p))
	})

	events := make(chan string, 4096)
	var monitors []*websocket.Conn

	http.Handle(p.URL, websocket.Handler(func(ws *websocket.Conn) {
		i := len(monitors)
		monitors = append(monitors, ws)
		defer func() {
			if e := ws.Close(); e != nil {
				log.Println(e.Error())
			}
			monitors[i] = nil
		}()

		var b struct{ Message string }
		for {
			if e := websocket.JSON.Receive(ws, &b); e == nil {
				events <- b.Message
			} else {
				log.Printf("socket receive error: %v\n", e)
				break
			}
		}
	}))

	go func() {
		for {
			e := <-events
			for _, ws := range monitors {
				if ws != nil {
					go func(ws *websocket.Conn, e string) {
						websocket.JSON.Send(ws, []interface{}{e})
					}(ws, e)
				}
			}
		}
	}()

	Abort(LAUNCH_FAILED, http.ListenAndServe(ADDRESS, nil))
}

func Abort(n int, e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(n)
	}
}

func BaseName() string {
	s := strings.Split(os.Args[0], "/")
	return s[len(s)-1]
}
