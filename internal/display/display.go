package display

import (
	"fmt"
	"strings"

	"github.com/ahmaruff/go-fleet/internal/game"
)

const (
	Reset  string = "\033[0m"
	Blue   string = "\033[34m"
	Green  string = "\033[32m"
	Red    string = "\033[31m"
	Yellow string = "\033[33m"
)

func ClearScreen() {
	fmt.Print("\033[2J\033[H")
}

func ConvertCellToChar(cellValue int) string {
	switch cellValue {
	case 0:
		return Blue + "~" + Reset // Water
	case 1:
		return Green + "S" + Reset // Ship (not hit)
	case 2:
		return Yellow + "O" + Reset // Miss
	case 3:
		return Red + "X" + Reset // Hit
	default:
		return "?" // Unknown
	}
}

func RenderGame(g *game.Game) {
	ClearScreen()

	isMyTurn := (g.CurrPlayer == 1)
	turnText := "Opponent's Turn"
	if isMyTurn {

		turnText = "Your Turn"
	}

	// Game header
	fmt.Println("===================================== GO-FLEET ==============================")

	if g.Phase == "PLAYING" {
		fmt.Printf("Player: %s vs %s | Phase: %s | Current Turn: %s\n",
			g.Player1.Name, g.Player2.Name, g.Phase, turnText)
	} else {
		fmt.Printf("Player: %s vs %s | Phase: %s\n",
			g.Player1.Name, g.Player2.Name, g.Phase)
	}

	fmt.Printf("Player: %s vs %s | Phase: %s | Current Turn: %s\n",
		g.Player1.Name, g.Player2.Name, g.Phase, turnText)
	fmt.Println("=============================================================================")

	// legends
	fmt.Printf(Blue + "~" + Reset + " = Water | " + Green + "S" + Reset + " = Ship | " + Red + "X" + Reset + " = Hit | " + Yellow + "O" + Reset + " = Miss\n")

	fmt.Println()

	fmt.Printf("Your Remaining Ships: %d\n", g.Player1.Board.ShipCount)
	fmt.Printf("Opponent's Remaining Ships: %d\n", g.Player2.Board.ShipCount)

	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Println()

	// Board headers
	fmt.Println("Your Board:                               Opponent's Board:")
	fmt.Println("   A B C D E F G H I J                        A B C D E F G H I J")

	// Render both boards side by side
	for row := 0; row < 10; row++ {
		// Your board (left side)
		fmt.Printf("%2d", row+1)
		for col := 0; col < 10; col++ {
			char := ConvertCellToChar(g.Player1.Board.Grid[row][col])
			fmt.Printf(" %s", char) // Space before each character
		}

		// Spacing between boards
		fmt.Print("                     ")

		// Opponent board (right side)
		fmt.Printf("%2d", row+1)
		for col := 0; col < 10; col++ {
			cellValue := g.Player2.Board.Grid[row][col]
			if cellValue == 2 || cellValue == 3 { // only render when HIT/MISS
				char := ConvertCellToChar(cellValue)
				fmt.Printf(" %s", char)
			} else {
				fmt.Printf(Blue + " ~" + Reset) // Consistent spacing
			}
		}
		fmt.Println()
	}

	fmt.Println()

	if g.Phase == "PLACING" {
		fmt.Println("Command: /set A1 — place your ship at A1")
	}

	if g.Phase == "PLAYING" {
		fmt.Println("Command: /fire B2 — fire at B2")
	}
}

func RenderGameAsString(g *game.Game) string {
	isMyTurn := (g.CurrPlayer == 1)
	turnText := "Opponent's Turn"
	if isMyTurn {

		turnText = "Your Turn"
	}

	var output strings.Builder

	// Game header
	output.WriteString("============================== GO-FLEET ==============================\n")

	if g.Phase == "PLAYING" {
		output.WriteString(fmt.Sprintf("Player: %s vs %s | Phase: %s | Current Turn: %s\n",
			g.Player1.Name, g.Player2.Name, g.Phase, turnText))
	} else {
		output.WriteString(fmt.Sprintf("Player: %s vs %s | Phase: %s\n",
			g.Player1.Name, g.Player2.Name, g.Phase))
	}

	output.WriteString("======================================================================\n")

	// legends
	output.WriteString(Blue + "~" + Reset + " = Water | " + Green + "S" + Reset + " = Ship | " + Red + "X" + Reset + " = Hit | " + Yellow + "O" + Reset + " = Miss\n")

	output.WriteString("\n")

	output.WriteString(fmt.Sprintf("Your Remaining Ships: %d\n", g.Player1.Board.ShipCount))
	output.WriteString(fmt.Sprintf("Opponent's Remaining Ships: %d\n", g.Player2.Board.ShipCount))

	output.WriteString("----------------------------------------------------------------------\n\n")

	// Board headers
	output.WriteString("Your Board:                               Opponent's Board:\n")
	output.WriteString("   A B C D E F G H I J                        A B C D E F G H I J\n")

	// Render both boards side by side
	for row := 0; row < 10; row++ {
		// Your board (left side)
		output.WriteString(fmt.Sprintf("%2d", row+1))
		for col := 0; col < 10; col++ {
			char := ConvertCellToChar(g.Player1.Board.Grid[row][col])
			output.WriteString(fmt.Sprintf(" %s", char)) // Space before each character
		}

		// Spacing between boards
		output.WriteString("                     ")

		// Opponent board (right side)
		output.WriteString(fmt.Sprintf("%2d", row+1))
		for col := 0; col < 10; col++ {
			cellValue := g.Player2.Board.Grid[row][col]
			if cellValue == 2 || cellValue == 3 { // only render when HIT/MISS
				char := ConvertCellToChar(cellValue)

				output.WriteString(fmt.Sprintf(" %s", char)) // Space before each character
			} else {
				output.WriteString(Blue + " ~" + Reset) // Space before each character
			}
		}
		output.WriteString("\n")
	}

	output.WriteString("\n")

	if g.Phase == "PLACING" {
		output.WriteString("Command: /set A1 — place your ship at A1\n")
	}

	if g.Phase == "PLAYING" {
		output.WriteString("Command: /fire B2 — fire at B2\n")
	}

	return output.String()
}
