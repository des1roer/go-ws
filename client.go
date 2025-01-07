package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

const (
	url     = "ws://localhost:8080/ws"
	message = "Hello from client!"
)

func main() {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(c)

	err = c.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Fatal(err)
	}

	_, p, err := c.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Received from server:", string(p))
}
