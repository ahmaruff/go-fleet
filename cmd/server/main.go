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
var games = make(map[*game.Game][2]net.Conn)
var waitingPlayer net.Conn

func main() {
	// Command line flag for port
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	fmt.Printf("[SERVER] Starting Go-Fleet Server on port %s...\n", *port)

	// Listen on specified port
	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatal("[SERVER] Failed to start server:", err)
	}

	defer listener.Close()

	fmt.Printf("[SERVER] Server listening on :%s\n", *port)

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			log.Println("[SERVER] Failed to accept connection:", err)
			continue
		}

		fmt.Println("[SERVER] New client connected!")

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
			fmt.Println("[SERVER] Client disconnected")
			return
		}

		message := strings.TrimSpace(string(buffer[:n]))
		fmt.Printf("[SERVER] Received: %s\n", message)

		response := handleCommand(conn, message) // Pass conn to track which client

		conn.Write([]byte(response + "\n"))
	}
}

func handleCommand(conn net.Conn, command string) string {
	parts := strings.Split(command, " ")

	switch parts[0] {
	case "/name":
		if len(parts) < 2 {
			return "[ERROR] - Usage: /name YourName"
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

		return "[NAME_SET] - Welcome " + playerName + "!"
	case "/ready":
		player := players[conn]

		if player == nil {
			return "[ERROR] - Please set your name first with /name"
		}

		if waitingPlayer == nil {
			// First player waiting
			waitingPlayer = conn
			return "[WAITING] - Looking for opponent..."
		} else {
			p1 := players[waitingPlayer]
			p2 := players[conn]

			newGame := game.Game{
				Player1:    p1,
				Player2:    p2,
				CurrPlayer: 1,
				Phase:      "PLACING",
			}

			games[&newGame] = [2]net.Conn{waitingPlayer, conn}

			// Notify both players
			waitingPlayer.Write([]byte("[GAME_START] - Match found! vs " + p2.Name + "\n"))
			conn.Write([]byte("[GAME_START] - Match found! vs " + p1.Name + "\n"))

			// Reset waiting player
			waitingPlayer = nil

			return ""
		}

	case "/set":
		if len(parts) < 2 {
			return "[ERROR] - Usage: /set A1"
		}

		coordinate := parts[1]

		return "[SHIP_PLACED] - Ship placed at " + coordinate
	case "/fire":
		if len(parts) < 2 {
			return "[ERROR] - Usage: /fire A1"
		}

		coordinate := parts[1]

		return "[SHOT_RESULT] - Fired at " + coordinate + " - MISS"
	default:
		return "[ERROR] - Unknown command"
	}
}
