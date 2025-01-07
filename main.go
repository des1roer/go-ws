package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Content string `json:"content"`
}

type ClientMessage struct {
	ClientId string  `json:"clientId"`
	Message  Message `json:"message"`
}

var clients = make(map[string]*websocket.Conn)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Позволяет всем запросам установить соединение
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			delete(clients, conn.RemoteAddr().String())
			break
		}

		switch messageType {
		case websocket.TextMessage:
			handleTextMessage(conn, string(p))
		case websocket.PongMessage:
			// Обработка PONG сообщений
		default:
			log.Printf("Unknown message type: %d", messageType)
		}
	}
}

func handleTextMessage(conn *websocket.Conn, message string) {
	var cm ClientMessage
	err := json.Unmarshal([]byte(message), &cm)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Received message from client %s: %s", conn.RemoteAddr(), message)
	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"timestamp":`+time.Now().Format(time.RFC3339)+`, "message": "`+message+`", "ClientId": "`+cm.ClientId+`"}`))
	if err != nil {
		log.Println(err)
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
