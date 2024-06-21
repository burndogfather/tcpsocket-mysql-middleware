package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

var (
	clientID  int
	mutex     = &sync.Mutex{}
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	mutex.Lock()
	clientID++
	id := clientID
	mutex.Unlock()

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("Client %d connected: %s\n", id, clientAddr)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading from client %d %s: %v\n", id, clientAddr, err)
			break
		}
		text = text[:len(text)-1] // Remove the newline character
		fmt.Printf("Received from client %d %s: %s\n", id, clientAddr, text)

		// Echo back the message to the same client
		_, err = writer.WriteString(fmt.Sprintf("Echo from server to client %d: %s\n", id, text))
		if err != nil {
			fmt.Printf("Error writing to client %d %s: %v\n", id, clientAddr, err)
			break
		}
		writer.Flush()
	}

	fmt.Printf("Client %d disconnected: %s\n", id, clientAddr)
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