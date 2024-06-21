package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

var (
	clients  = make(map[net.Conn]bool)
	messages = make(chan string)
	mutex    = &sync.Mutex{}
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("Client connected: %s\n", clientAddr)

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	reader := bufio.NewReader(conn)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Printf("Received from %s: %s", clientAddr, text)
		messages <- fmt.Sprintf("%s: %s", clientAddr, text)
	}

	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()

	fmt.Printf("Client disconnected: %s\n", clientAddr)
}

func broadcastMessages() {
	for {
		msg := <-messages
		fmt.Println("Broadcasting:", msg)
		mutex.Lock()
		for conn := range clients {
			fmt.Fprintln(conn, msg)
		}
		mutex.Unlock()
	}
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

	go broadcastMessages()
	go acceptConnections(listener)

	// Handle server input
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) > 0 {
			messages <- fmt.Sprintf("Server: %s", text)
		}
	}
}