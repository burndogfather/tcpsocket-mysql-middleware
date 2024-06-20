package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"github.com/pires/go-proxyproto"
)

func main() {
	// TCP 서버를 시작합니다.
	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	// PROXY 프로토콜을 지원하는 리스너로 감쌉니다.
	proxyListener := &proxyproto.Listener{Listener: listener}
	defer proxyListener.Close()

	log.Println("Server is listening on port 4000")

	for {
		// 클라이언트의 연결을 수락합니다.
		conn, err := proxyListener.Accept()
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
	remoteAddr := conn.RemoteAddr().String()

	log.Printf("Client connected: %s", remoteAddr)

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

	log.Printf("Client disconnected: %s", remoteAddr)
}