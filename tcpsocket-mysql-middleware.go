package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

var (
	mutex = &sync.Mutex{}
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("Client connected: %s\n", clientAddr)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		text, err := reader.ReadString('::')
		if err != nil {
			fmt.Printf("Error reading from client %s: %v\n", clientAddr, err)
			break
		}
		text = text[:len(text)-1] // Remove the newline character
		fmt.Printf("Received from %s: %s\n", clientAddr, text)
		
		// Echo back the message to the same client
		_, err = writer.WriteString(fmt.Sprintf("Echo from server: %s\n", text))
		if err != nil {
			fmt.Printf("Error writing to client %s: %v\n", clientAddr, err)
			break
		}
		writer.Flush()
	}

	fmt.Printf("Client disconnected: %s\n", clientAddr)
}

func acceptConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started on port 4000")

	go acceptConnections(listener)

	// Handle server input (optional)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) > 0 {
			fmt.Printf("Server input: %s\n", text)
		}
	}
}