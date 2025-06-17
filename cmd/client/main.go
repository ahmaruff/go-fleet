package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	// Command line flags
	host := flag.String("host", "localhost", "Server host")
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	address := *host + ":" + *port
	fmt.Printf("[INFO] - Connecting to Go-Fleet Server at %s...\n", address)

	// Connect to server
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal("[ERROR] - Failed to connect to server:", err)
	}
	defer conn.Close()

	fmt.Println("[INFO] - Connected!")

	// Ask for player name
	fmt.Print(">> Please enter your name: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	playerName := scanner.Text()

	// Send name to server
	_, err = conn.Write([]byte("/name " + playerName + "\n"))
	if err != nil {
		log.Fatal("[ERROR] - Failed to send name:", err)
	}

	// Start listening for server messages
	go listenForMessages(conn)

	fmt.Println(">> Type commands (/ready, /set A1, /fire B2) or 'quit' to exit:")

	// Continue with existing input loop...
	for scanner.Scan() {
		message := scanner.Text()

		if message == "quit" {
			break
		}

		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			log.Println("[ERROR] - Failed to send message:", err)
			break
		}
	}
}

func listenForMessages(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	var displayBuffer strings.Builder

	inDisplayMode := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "DISPLAY_UPDATE" {
			inDisplayMode = true

			displayBuffer.Reset()
			continue
		}

		if line == "END_DISPLAY" {
			inDisplayMode = false

			// Clear screen and show the game board
			fmt.Print("\033[2J\033[H") // Clear screen
			fmt.Print(displayBuffer.String())
			fmt.Print("\nEnter command: ")
			continue
		}

		if inDisplayMode {
			displayBuffer.WriteString(line + "\n")
			continue
		}

		// Regular server messages
		fmt.Printf("\n# %s\n", line)
		fmt.Print("Enter command: ")
	}
}
