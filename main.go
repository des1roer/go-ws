package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
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
			var cm ClientMessage
			err := json.Unmarshal(p, &cm)
			if err == nil {
				clients[cm.ClientId] = conn
				handleTextMessage(conn, cm.ClientId, cm.Message.Content)
			} else {
				log.Printf("Error unmarshaling message from client %s", conn.RemoteAddr())
			}
		case websocket.PongMessage:
			// Обработка PONG сообщений
		default:
			log.Printf("Unknown message type: %d", messageType)
		}
	}
}

func handleTextMessage(conn *websocket.Conn, clientId string, message string) {
	log.Printf("Received message from client %s: %s", clientId, message)

	// Проверяем наличие клиента в списке соединений
	if _, ok := clients[clientId]; ok {
		err := conn.WriteMessage(websocket.TextMessage, []byte(`{"timestamp":`+time.Now().Format(time.RFC3339)+`, "message": "`+message+`"}"`))
		if err != nil {
			log.Println(err)
			delete(clients, clientId)
		}

		log.Printf("Sent hello message to client %s", clientId)
	} else {
		log.Printf("Client %s not found in connections", clientId)
	}
}

func getApiClientStatus(w http.ResponseWriter, r *http.Request) {
	clientId := r.URL.Query().Get("clientId")

	if clientId == "" {
		http.Error(w, "Missing clientId parameter", http.StatusBadRequest)
		return
	}

	if _, ok := clients[clientId]; ok {
		w.WriteHeader(http.StatusOK)
		handleTextMessage(clients[clientId], clientId, `HI`)
		_, err := w.Write([]byte(`{"status":"connected"}`))
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"status":"not connected"}`))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/api/status", getApiClientStatus)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
