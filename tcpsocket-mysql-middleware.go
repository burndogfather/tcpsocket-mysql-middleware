package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		// 메시지를 읽습니다.
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Received: %s", message)

		// 메시지를 다시 보냅니다.
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("WebSocket server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}