package display

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ahmaruff/go-fleet/internal/game"
)

func ClearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls") // Windows
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear") // Unix-based systems
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func ConvertCellToChar(cellValue int) rune {
	switch cellValue {
	case 0:
		return '~' // Empty water
	case 1:
		return 'S' // Ship (not hit)
	case 2:
		return 'O' // Miss
	case 3:
		return 'X' // Hit
	default:
		return '?' // Unknown
	}
}

func RenderGame(g *game.Game) {
	ClearScreen()

	// Game header
	fmt.Println("========================================== GO-FLEET ==========================================")
	fmt.Printf("Player: %s vs %s | Phase: %s | Current Turn: Player %d\n",
		g.Player1.Name, g.Player2.Name, g.Phase, g.CurrPlayer)
	fmt.Println("===========================================================================================")

	// Board headers
	fmt.Println("Your Board:                               Opponent's Board:")
	fmt.Println("   A B C D E F G H I J                        A B C D E F G H I J")

	// Render both boards side by side
	for row := 0; row < 10; row++ {
		// Your board (left side)
		fmt.Printf("%2d", row+1)
		for col := 0; col < 10; col++ {
			char := ConvertCellToChar(g.Player1.Board.Grid[row][col])
			fmt.Printf(" %c", char) // Space before each character
		}

		// Spacing between boards
		fmt.Print("                     ")

		// Opponent board (right side)
		fmt.Printf("%2d", row+1)
		for col := 0; col < 10; col++ {
			cellValue := g.Player2.Board.Grid[row][col]
			if cellValue == 2 || cellValue == 3 { // only render when HIT/MISS
				char := ConvertCellToChar(cellValue)
				fmt.Printf(" %c", char)
			} else {
				fmt.Printf(" ~") // Consistent spacing
			}
		}
		fmt.Println()
	}
}

func RenderGameAsString(g *game.Game) string {
	var output strings.Builder

	// Game header
	output.WriteString("========================================== GO-FLEET ==========================================")
	output.WriteString(fmt.Sprintf("Player: %s vs %s | Phase: %s | Current Turn: Player %d\n",
		g.Player1.Name, g.Player2.Name, g.Phase, g.CurrPlayer))
	output.WriteString("===========================================================================================")

	// Board headers
	output.WriteString("Your Board:                               Opponent's Board:")
	output.WriteString("   A B C D E F G H I J                        A B C D E F G H I J")

	// Render both boards side by side
	for row := 0; row < 10; row++ {
		// Your board (left side)
		output.WriteString(fmt.Sprintf("%2d", row+1))
		for col := 0; col < 10; col++ {
			char := ConvertCellToChar(g.Player1.Board.Grid[row][col])
			output.WriteString(fmt.Sprintf(" %c", char)) // Space before each character
		}

		// Spacing between boards
		output.WriteString("                     ")

		// Opponent board (right side)
		output.WriteString(fmt.Sprintf("%2d", row+1))
		for col := 0; col < 10; col++ {
			cellValue := g.Player2.Board.Grid[row][col]
			if cellValue == 2 || cellValue == 3 { // only render when HIT/MISS
				char := ConvertCellToChar(cellValue)

				output.WriteString(fmt.Sprintf(" %c", char)) // Space before each character
			} else {
				output.WriteString(" ~") // Space before each character
			}
		}
		output.WriteString("\n")
	}

	output.WriteString("\nCommands: /ready, /set A1, /fire B2\n")

	return output.String()
}
