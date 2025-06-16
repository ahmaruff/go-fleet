package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/ahmaruff/go-fleet/internal/display"
	"github.com/ahmaruff/go-fleet/internal/game"
)

var players = make(map[net.Conn]*game.Player)
var games = make(map[*game.Game][2]net.Conn)
var waitingPlayer net.Conn

func findGameByConnection(conn net.Conn) *game.Game {
	for gameInstance, connections := range games {
		if connections[0] == conn || connections[1] == conn {
			return gameInstance
		}
	}
	return nil
}

func getCurrentPlayerShips(g *game.Game, conn net.Conn) int {
	connections := games[g]
	if connections[0] == conn {
		return g.Player1.Board.ShipCount
	}
	return g.Player2.Board.ShipCount
}

// Function to capture display output for a specific player
func captureDisplayForPlayer(gameInstance *game.Game, playerConn net.Conn) string {
	// Find which player this connection represents
	connections := games[gameInstance]
	isPlayer1 := connections[0] == playerConn

	// Create a version of the game from this player's perspective
	var playerGame *game.Game
	if isPlayer1 {
		// Player 1's perspective - they are "Player1" in the game
		playerGame = gameInstance
	} else {
		// Player 2's perspective - swap players so they see themselves as "Player1"
		playerGame = &game.Game{
			Player1:    gameInstance.Player2, // They see themselves as Player1
			Player2:    gameInstance.Player1, // Opponent as Player2
			CurrPlayer: gameInstance.CurrPlayer,
			Phase:      gameInstance.Phase,
		}
	}

	return display.RenderGameAsString(playerGame)
}

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

		currentGame := findGameByConnection(conn)

		if currentGame == nil {
			return "[ERROR] - You're not in a game. Use /ready first"
		}

		if currentGame.Phase != "PLACING" {
			return "[ERROR] - Not in placement phase"
		}

		player := players[conn]
		coordinate := parts[1]
		success := currentGame.PlaceShipForPlayer(player, coordinate)

		if !success {
			return "[ERROR] - Cannot place ship at " + coordinate
		}

		response := fmt.Sprintf("[SHIP_PLACED] - Ship placed at %s (%d/5)", coordinate, getCurrentPlayerShips(currentGame, conn))

		if currentGame.Phase == "PLAYING" {
			// Both players have 5 ships, game started!
			connections := games[currentGame]
			connections[0].Write([]byte("[COMBAT_START] - All ships placed! Combat phase begins!\n"))
			connections[1].Write([]byte("[COMBAT_START] - All ships placed! Combat phase begins!\n"))
		}

		return response
	case "/fire":
		if len(parts) < 2 {
			return "[ERROR] - Usage: /fire A1"
		}

		currentGame := findGameByConnection(conn)
		if currentGame == nil {
			return "[ERROR] - You're not in a game"
		}

		if currentGame.Phase != "PLAYING" {
			return "[ERROR] - Not in combat phase"
		}

		player := players[conn]
		coordinate := parts[1]

		result := currentGame.FireAtOpponent(player, coordinate)

		fireMsg := map[int]string{
			3: "HIT",
			2: "MISS",
		}

		resultMsg := fireMsg[result]

		response := fmt.Sprintf("[SHOT_RESULT] - %s at %s", resultMsg, coordinate)

		// Check if game is over
		winner, gameOver := currentGame.IsGameOver()
		if gameOver {
			connections := games[currentGame]
			winnerName := ""
			if winner == 1 {
				winnerName = currentGame.Player1.Name
			} else {
				winnerName = currentGame.Player2.Name
			}

			connections[0].Write([]byte("======================================\n"))
			connections[0].Write([]byte("[GAME_OVER] - " + winnerName + " wins!\n"))
			connections[0].Write([]byte("======================================\n"))

			connections[1].Write([]byte("======================================\n"))
			connections[1].Write([]byte("[GAME_OVER] - " + winnerName + " wins!\n"))
			connections[1].Write([]byte("======================================\n"))
		}

		return response

	default:
		return "[ERROR] - Unknown command"
	}
}
