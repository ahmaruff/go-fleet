package main

import (
	"flag"
	"fmt"
	"github.com/ahmaruff/go-fleet/internal/game"
	"log"
	"net"
	"strings"
)

var players = make(map[net.Conn]*game.Player)

func main() {
	// Command line flag for port
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	fmt.Printf("Starting Go-Fleet Server on port %s...\n", *port)

	// Listen on specified port
	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}

	defer listener.Close()

	fmt.Printf("Server listening on :%s\n", *port)

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}

		fmt.Println("New client connected!")

		// Handle each client in a separate goroutine
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	defer delete(players, conn) // Clean up when client disconnects

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Client disconnected")
			return
		}

		message := strings.TrimSpace(string(buffer[:n]))
		fmt.Printf("Received: %s\n", message)

		response := handleCommand(conn, message) // Pass conn to track which client

		conn.Write([]byte(response + "\n"))
	}
}

func handleCommand(conn net.Conn, command string) string {
	parts := strings.Split(command, " ")

	switch parts[0] {
	case "/name":
		if len(parts) < 2 {
			return "ERROR - Usage: /name YourName"
		}
		playerName := strings.Join(parts[1:], " ")

		// Create Player object
		board := &game.Board{ShipCount: 0}
		player := &game.Player{
			Name:  playerName,
			Board: board,
		}

		// Store player for this connection
		players[conn] = player

		return "NAME_SET - Welcome " + playerName + "!"
	case "/ready":
		return "READY_ACK - Waiting for opponent..."
	case "/set":
		if len(parts) < 2 {
			return "ERROR - Usage: /set A1"
		}

		coordinate := parts[1]

		return "SHIP_PLACED - Ship placed at " + coordinate
	case "/fire":
		if len(parts) < 2 {
			return "ERROR - Usage: /fire A1"
		}

		coordinate := parts[1]

		return "SHOT_RESULT - Fired at " + coordinate + " - MISS"
	default:
		return "ERROR - Unknown command"
	}
}
