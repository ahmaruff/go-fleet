package main

import (
	"fmt"

	"github.com/ahmaruff/go-fleet/internal/game"
)

func main() {
	b1 := game.Board{
		Grid:      [10][10]int{},
		ShipCount: 0,
	}

	p1 := game.Player{
		Name:  "Player 1",
		Board: &b1,
	}

	b2 := game.Board{
		Grid:      [10][10]int{},
		ShipCount: 0,
	}

	p2 := game.Player{
		Name:  "Player 2",
		Board: &b2,
	}

	g := game.NewGame(&p1, &p2)

	// TEST 1: Check initial game state
	fmt.Println("=== INITIAL STATE ===")
	fmt.Println("Current Player:", g.CurrPlayer)
	fmt.Println("Phase:", g.Phase)
	fmt.Println("P1 Ships:", g.Player1.Board.ShipCount)
	fmt.Println("P2 Ships:", g.Player2.Board.ShipCount)

	// TEST 2: Try placing ships
	fmt.Println("\n=== PLACING SHIPS ===")
	resTest2 := g.PlaceShipForPlayer("A1")
	fmt.Println("Placed ship at A1:", resTest2)
	fmt.Println("P1 Ships after:", g.Player1.Board.ShipCount)

	// TEST 3: Place more ships and test player switching
	fmt.Println("\n=== MORE SHIP PLACEMENT ===")
	g.PlaceShipForPlayer("B1")
	g.PlaceShipForPlayer("C1")
	g.PlaceShipForPlayer("D1")
	g.PlaceShipForPlayer("E1")
	fmt.Println("P1 Ships after 5 placements:", g.Player1.Board.ShipCount)

	// Switch to player 2
	g.SwitchPlayer()
	fmt.Println("Current Player after switch:", g.CurrPlayer)

	// Place ships for player 2
	g.PlaceShipForPlayer("A2")
	g.PlaceShipForPlayer("B2")
	g.PlaceShipForPlayer("C2")
	g.PlaceShipForPlayer("D2")
	g.PlaceShipForPlayer("E2")
	fmt.Println("P2 Ships after placement:", g.Player2.Board.ShipCount)
	fmt.Println("Phase after both players have 5 ships:", g.Phase)

	// TEST 4: Combat phase
	fmt.Println("\n=== COMBAT PHASE ===")
	fmt.Println("Current Phase:", g.Phase)

	// Player 2 fires at Player 1's ships
	resTest41 := g.FireAtOpponent("A1")                      // Should hit P1's ship
	fmt.Println("Player 2 fires at A1 - Result:", resTest41) // Should be 3 (hit)
	fmt.Println("P1 Ships remaining:", g.Player1.Board.ShipCount)

	// Test a miss
	resTest42 := g.FireAtOpponent("A10")                      // Should miss (empty water)
	fmt.Println("Player 2 fires at A10 - Result:", resTest42) // Should be 2 (miss)

	// Test game over scenario
	fmt.Println("\n=== TESTING GAME OVER ===")
	// Fire at all P1's remaining ships
	g.FireAtOpponent("B1")
	g.FireAtOpponent("C1")
	g.FireAtOpponent("D1")
	g.FireAtOpponent("E1")
	fmt.Println("P1 Ships after all hits:", g.Player1.Board.ShipCount)
	fmt.Println("Final Phase:", g.Phase)

	winner, gameOver := g.IsGameOver()
	fmt.Println("Game Over:", gameOver, "Winner:", winner)
}
