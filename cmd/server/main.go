package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/ahmaruff/go-fleet/internal/display"
	"github.com/ahmaruff/go-fleet/internal/effects"
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

		swappedCurrPlayer := gameInstance.CurrPlayer
		if swappedCurrPlayer == 1 {
			swappedCurrPlayer = 2
		} else {
			swappedCurrPlayer = 1
		}

		playerGame = &game.Game{
			Player1:    gameInstance.Player2, // They see themselves as Player1
			Player2:    gameInstance.Player1, // Opponent as Player2
			CurrPlayer: swappedCurrPlayer,
			Phase:      gameInstance.Phase,
		}
	}

	return display.RenderGameAsString(playerGame)
}

func main() {
	display.ClearScreen()

	welcomeEffect := effects.GetEffect("WELCOME")

	fmt.Println()
	fmt.Printf("%s", welcomeEffect)
	fmt.Println()
	fmt.Println()

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

			currentGame := findGameByConnection(conn)
			if currentGame != nil {
				// Get connections BEFORE deleting
				connections := games[currentGame]

				// Remove game from tracking
				delete(games, currentGame)

				// Notify the remaining player (if still connected)
				for _, connection := range connections {
					if connection != conn { // Don't send to disconnected player
						connection.Write([]byte("OPPONENT_DISCONNECTED\n"))
						connection.Write([]byte("GAME_RESET\n"))
					}
				}
			}

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

		if _, exists := players[conn]; exists {
			return "[ERROR] - You already have a name set. You can't change it."
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

		welcomeEffect := effects.GetEffect("WELCOME")

		conn.Write([]byte("EFFECT_UPDATE\n" + welcomeEffect + "\nEFFECT_END\n"))

		return "[NAME_SET] - Welcome " + playerName + "!"
	case "/ready":
		player := players[conn]

		if player == nil {
			return "[ERROR] - Please set your name first with /name"
		}

		if waitingPlayer == nil {
			// First player waiting
			waitingPlayer = conn
			// Send waiting effect
			effectOutput := effects.GetEffect("WAITING")
			conn.Write([]byte("\nEFFECT_UPDATE\n" + effectOutput + "\nEFFECT_END\n"))
			return "[WAITING] - Looking for opponent..."
		}

		if conn == waitingPlayer {

			// Send waiting effect
			effectOutput := effects.GetEffect("WAITING")
			conn.Write([]byte("EFFECT_UPDATE\n" + effectOutput + "\nEFFECT_END\n"))

			return "[WAITING] - Looking for opponent..."
		}

		currentGame := findGameByConnection(conn)

		if currentGame != nil {
			return "[ERROR] - You are already in game, unable to use /ready command, use /set or /fire"
		}

		p1 := players[waitingPlayer]
		p2 := players[conn]

		newGame := game.NewGame(p1, p2)
		games[newGame] = [2]net.Conn{waitingPlayer, conn}

		// Notify both players
		waitingPlayer.Write([]byte("[GAME_START] - Match found! vs " + p2.Name + "\n"))
		conn.Write([]byte("[GAME_START] - Match found! vs " + p1.Name + "\n"))

		// effect match found
		matchEffect := effects.GetEffect("MATCH_FOUND")
		waitingPlayer.Write([]byte("EFFECT_UPDATE\n" + matchEffect + "\nEFFECT_END\n"))
		conn.Write([]byte("EFFECT_UPDATE\n" + matchEffect + "\nEFFECT_END\n"))

		display1 := captureDisplayForPlayer(newGame, waitingPlayer)
		display2 := captureDisplayForPlayer(newGame, conn)

		waitingPlayer.Write([]byte("DISPLAY_UPDATE\n" + display1 + "END_DISPLAY\n"))
		conn.Write([]byte("DISPLAY_UPDATE\n" + display2 + "END_DISPLAY\n"))

		// Reset waiting player
		waitingPlayer = nil

		return ""

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

		if player.Board.ShipCount > 4 {
			vesselReadyEffect := effects.GetEffect("ALL_SHIPS_READY")
			conn.Write([]byte("EFFECT_UPDATE\n" + vesselReadyEffect + "\nEFFECT_END\n"))
		}

		if currentGame.Phase == "PLAYING" {
			// Both players have 5 ships, game started!
			connections := games[currentGame]
			connections[0].Write([]byte("[COMBAT_START] - All ships placed! Combat phase begins!\n"))
			connections[1].Write([]byte("[COMBAT_START] - All ships placed! Combat phase begins!\n"))

			battleStartEffect := effects.GetEffect("BATTLE_START")
			connections[0].Write([]byte("EFFECT_UPDATE\n" + battleStartEffect + "\nEFFECT_END\n"))
			connections[1].Write([]byte("EFFECT_UPDATE\n" + battleStartEffect + "\nEFFECT_END\n"))

			// Send display update to both players when combat starts
			display1 := captureDisplayForPlayer(currentGame, connections[0])
			display2 := captureDisplayForPlayer(currentGame, connections[1])

			connections[0].Write([]byte("DISPLAY_UPDATE\n" + display1 + "END_DISPLAY\n"))
			connections[1].Write([]byte("DISPLAY_UPDATE\n" + display2 + "END_DISPLAY\n"))
		} else {
			// Normal ship placement - send display only to current player
			displayOutput := captureDisplayForPlayer(currentGame, conn)
			conn.Write([]byte("DISPLAY_UPDATE\n" + displayOutput + "END_DISPLAY\n"))
		}

		response := fmt.Sprintf("[SHIP_PLACED] - Ship placed at %s (%d/5)", coordinate, getCurrentPlayerShips(currentGame, conn))

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

		connections := games[currentGame]
		isPlayer1 := connections[0] == conn
		playerNumber := 1
		if !isPlayer1 {
			playerNumber = 2
		}

		if currentGame.CurrPlayer != playerNumber {
			return "[ERROR] - Not your turn! Wait for opponent to fire."
		}

		player := players[conn]
		coordinate := parts[1]

		result := currentGame.FireAtOpponent(player, coordinate)

		if result != 3 && result != 2 {
			conn.Write([]byte("[ERROR] - Invalid shot at " + coordinate + "\n"))
			return ""
		}

		fireMsg := map[int]string{
			3: "HIT",
			2: "MISS",
		}

		resultMsg := fireMsg[result]
		response := fmt.Sprintf("[SHOT_RESULT] - %s at %s", resultMsg, coordinate)

		fireEffect := effects.GetEffect(resultMsg)
		connections[0].Write([]byte("EFFECT_UPDATE\n" + fireEffect + "\nEFFECT_END\n"))
		connections[1].Write([]byte("EFFECT_UPDATE\n" + fireEffect + "\nEFFECT_END\n"))

		display1 := captureDisplayForPlayer(currentGame, connections[0])
		display2 := captureDisplayForPlayer(currentGame, connections[1])

		connections[0].Write([]byte("DISPLAY_UPDATE\n" + display1 + "END_DISPLAY\n"))
		connections[1].Write([]byte("DISPLAY_UPDATE\n" + display2 + "END_DISPLAY\n"))

		// Check if game is over
		winner, gameOver := currentGame.IsGameOver()

		defeatIndex := 0
		winnerIndex := 1
		if gameOver {
			connections := games[currentGame]
			winnerName := ""
			if winner == 1 {
				winnerName = currentGame.Player1.Name
				winnerIndex = 0
				defeatIndex = 1
			} else {
				winnerName = currentGame.Player2.Name

				winnerIndex = 1
				defeatIndex = 0
			}

			victoryEffect := effects.GetEffect("VICTORY")
			defeatEffect := effects.GetEffect("DEFEAT")

			connections[winnerIndex].Write([]byte("EFFECT_UPDATE\n" + victoryEffect + "\nEFFECT_END\n"))
			connections[defeatIndex].Write([]byte("EFFECT_UPDATE\n" + defeatEffect + "\nEFFECT_END\n"))

			connections[0].Write([]byte("======================================\n"))
			connections[0].Write([]byte("[GAME_OVER] - " + winnerName + " wins!\n"))
			connections[0].Write([]byte("======================================\n"))

			connections[1].Write([]byte("======================================\n"))
			connections[1].Write([]byte("[GAME_OVER] - " + winnerName + " wins!\n"))
			connections[1].Write([]byte("======================================\n"))

			// CLEANUP: Remove game from tracking
			delete(games, currentGame)

			// Send reset messages to both players
			connections[0].Write([]byte("GAME_RESET\n"))
			connections[1].Write([]byte("GAME_RESET\n"))
		}

		return response

	default:
		return "[ERROR] - Unknown command"
	}
}
