package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

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
	// Здесь можно добавить логику обработки входящих сообщений
	log.Printf("Received message: %s", message)
	// Отправка ответного сообщения
	err := conn.WriteMessage(websocket.TextMessage, []byte("Server received your message"))
	if err != nil {
		log.Println(err)
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
