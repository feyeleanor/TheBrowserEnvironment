package main

import (
	"log"

	"golang.org/x/net/websocket"
)

const SERVER = "ws://localhost:3000/socket"
const ADDRESS = "http://localhost/"

func main() {
	if ws, e := websocket.Dial(SERVER, "", ADDRESS); e == nil {
		var m []interface{}
		for {
			if e := websocket.JSON.Receive(ws, &m); e == nil {
				log.Printf("message: %v\n", m)
			} else {
				log.Fatal(e)
			}
		}
	} else {
		log.Fatal(e)
	}
}
