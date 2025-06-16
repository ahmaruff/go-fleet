package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// Command line flags
	host := flag.String("host", "localhost", "Server host")
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	address := *host + ":" + *port
	fmt.Printf("Connecting to Go-Fleet Server at %s...\n", address)

	// Connect to server
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}
	defer conn.Close()

	fmt.Println("Connected!")

	// Ask for player name
	fmt.Print("Enter your name: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	playerName := scanner.Text()

	// Send name to server
	_, err = conn.Write([]byte("/name " + playerName + "\n"))
	if err != nil {
		log.Fatal("Failed to send name:", err)
	}

	// Start listening for server messages
	go listenForMessages(conn)

	fmt.Println("Type commands (/ready, /set A1, /fire B2) or 'quit' to exit:")

	// Continue with existing input loop...
	for scanner.Scan() {
		message := scanner.Text()

		if message == "quit" {
			break
		}

		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			log.Println("Failed to send message:", err)
			break
		}
	}
}

func listenForMessages(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Disconnected from server")
			return
		}

		fmt.Printf("Server: %s", string(buffer[:n]))
	}
}
