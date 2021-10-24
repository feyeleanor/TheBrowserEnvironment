package main

import (
	"log"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

const SERVER = "ws://localhost:3000/socket"
const ADDRESS = "http://localhost/"

func main() {
	if ws, e := websocket.Dial(SERVER, "", ADDRESS); e == nil {
		var b struct{ Message string }
		for _, v := range os.Args[1:] {
			b.Message = v
			websocket.JSON.Send(ws, &b)
			time.Sleep(1 * time.Second)
		}
	} else {
		log.Fatal(e)
	}
}
