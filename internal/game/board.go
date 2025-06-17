package game

// GRID CELL STATE ----

// 0 = Empty water (not fired upon)
// 1 = Ship present (not fired upon)
// 2 = Water hit (miss)
// 3 = Ship hit

type Board struct {
	Grid      [10][10]int
	ShipCount int
}

func (b *Board) PlaceShip(row, col int) bool {
	if row < 0 || col < 0 || row > 9 || col > 9 {
		return false // Invalid coordinates
	}

	if b.ShipCount >= 5 {
		return false
	}

	if b.Grid[row][col] != 0 {
		return false // Cell already has something
	}

	b.Grid[row][col] = 1
	b.ShipCount += 1

	return true
}

func (b *Board) Fire(row, col int) int {
	fireMap := map[int]int{
		0: 2,
		1: 3,
		2: 2,
		3: 3,
	}

	currentStatus := b.Grid[row][col]
	nextStatus := fireMap[currentStatus]
	b.Grid[row][col] = nextStatus

	if nextStatus == 3 {
		b.ShipCount--
	}
	return nextStatus
}

func (b *Board) IsValidPosition(row, col int) bool {
	if row >= 0 && row <= 9 && col >= 0 && col <= 9 {
		return true
	}

	return false
}

func (b *Board) AllShipDestroyed() bool {
	if b.ShipCount <= 0 {
		return true
	}

	return false
}
