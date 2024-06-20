package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	// TCP 서버를 시작합니다.
	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()
	log.Println("Server is listening on port 4000")

	for {
		// 클라이언트의 연결을 수락합니다.
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		// 클라이언트 연결을 처리합니다.
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Client connected: %s", conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)
	for {
		// 클라이언트로부터 메시지를 읽습니다.
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			break
		}

		message = strings.TrimSpace(message)
		log.Printf("Received message: %s", message)

		// 클라이언트에게 응답을 보냅니다.
		response := fmt.Sprintf("You said: %s\n", message)
		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Printf("Failed to send response: %v", err)
			break
		}
	}

	log.Printf("Client disconnected: %s", conn.RemoteAddr().String())
}