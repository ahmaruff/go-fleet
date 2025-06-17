package game

import (
	"fmt"
	"strings"
	"unicode"
)

func ConvertCell(cell string) (int, int) {
	cell = strings.ToUpper(cell)

	// example: cell = AB41;
	// so this is separating "AB" (column) and "41" (row)

	var letterPart, numberPart string
	for _, ch := range cell {
		if unicode.IsLetter(ch) {
			letterPart += string(ch)
		} else if unicode.IsDigit(ch) {
			numberPart += string(ch)
		}
	}

	col := 0
	for _, ch := range letterPart {
		ascii := int(ch)
		offset := ascii - 65 // A=0, B=1, ..., Z=25
		col = col*26 + offset + 1
	}

	// karena A = 0
	col--

	var row int
	fmt.Sscanf(numberPart, "%d", &row)

	// karena baris dimulai dari 1, sedangkan index array dari 0
	row--

	// Validate range
	if row < 0 || row > 9 || col < 0 || col > 9 {
		return -1, -1 // Invalid coordinate
	}

	return row, col
}
