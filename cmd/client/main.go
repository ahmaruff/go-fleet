package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/ahmaruff/go-fleet/internal/display"
	"github.com/ahmaruff/go-fleet/internal/effects"
)

func main() {
	display.ClearScreen()
	welcomeEffect := effects.GetEffect("WELCOME")

	fmt.Println()
	fmt.Printf("%s\n", welcomeEffect)
	fmt.Println()

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

	// Small delay to let server response come through
	time.Sleep(100 * time.Millisecond)

	display.ClearScreen()

	fmt.Println("============================== GO-FLEET ==============================")
	fmt.Println("Type '/ready' if you're ready for war or '/quit' to exit")
	fmt.Println("======================================================================")
	fmt.Println()

	// Continue with existing input loop...
	for scanner.Scan() {
		message := scanner.Text()

		if message == "quit" || message == "/quit" || message == "/exit" {
			break
		}
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			log.Println("[ERROR] - Failed to send message:", err)
			break
		}
	}
}

var effectQueue []string
var currentlyShowingEffect bool
var queuedDisplay string

func listenForMessages(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	var displayBuffer strings.Builder
	var effectBuffer strings.Builder

	inDisplayMode := false
	inEffectMode := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "\n" {
			continue
		}

		if line == "OPPONENT_DISCONNECTED" {
			// Clear all game state
			effectQueue = nil
			currentlyShowingEffect = false
			queuedDisplay = ""

			fmt.Println("---------------------")
			fmt.Println("Opponent Disconected!")
			fmt.Println("---------------------")

			time.Sleep(3 * time.Second)

			// Show ready prompt
			display.ClearScreen()

			fmt.Println("============================== GO-FLEET ==============================")
			fmt.Println("Type '/ready' if you're ready for war or '/quit' to exit")
			fmt.Println("======================================================================")
			fmt.Println()
			continue
		}

		if line == "GAME_RESET" {
			// Clear all game state
			effectQueue = nil
			currentlyShowingEffect = false
			queuedDisplay = ""

			// Show ready prompt
			display.ClearScreen()

			fmt.Println("============================== GO-FLEET ==============================")
			fmt.Println("Type '/ready' if you're ready for war or '/quit' to exit")
			fmt.Println("======================================================================")
			fmt.Println()
			continue
		}

		if line == "EFFECT_UPDATE" {
			//			fmt.Printf("[DEBUG] Starting effect mode\n")
			inEffectMode = true
			effectBuffer.Reset()
			continue
		}

		if line == "EFFECT_END" {
			inEffectMode = false
			effectQueue = append(effectQueue, effectBuffer.String())

			if !currentlyShowingEffect {
				showNextEffect()
			}
			continue
		}

		if inEffectMode {
			effectBuffer.WriteString(line + "\n")
			continue
		}

		if line == "DISPLAY_UPDATE" {
			inDisplayMode = true
			displayBuffer.Reset()
			continue
		}

		if line == "END_DISPLAY" {
			inDisplayMode = false
			queuedDisplay = displayBuffer.String()

			// If no effects are showing, display immediately
			if !currentlyShowingEffect {
				display.ClearScreen()
				fmt.Print(queuedDisplay)
				queuedDisplay = ""
			}
			continue
		}

		if inDisplayMode {
			displayBuffer.WriteString(line + "\n")
			continue
		}

		// Regular server messages
		if !currentlyShowingEffect {
			fmt.Printf("%s\n", line)
		}
	}
}

func showNextEffect() {
	if len(effectQueue) == 0 {
		currentlyShowingEffect = false
		// Show queued display if available
		if queuedDisplay != "" {
			display.ClearScreen()
			fmt.Print(queuedDisplay)
			queuedDisplay = ""
		}
		return
	}

	currentlyShowingEffect = true
	effect := effectQueue[0]
	effectQueue = effectQueue[1:] // Remove first effect

	display.ClearScreen()
	fmt.Print(effect)

	// Timer to show next effect
	go func() {
		time.Sleep(3 * time.Second)
		showNextEffect() // Recursively show next effect
	}()
}
