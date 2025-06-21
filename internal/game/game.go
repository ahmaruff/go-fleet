package game

// PHASE STATUS
// PLACING
// PLAYING
// FINISHED

type Game struct {
	Player1    *Player
	Player2    *Player
	CurrPlayer int
	Phase      string
}

func NewGame(p1, p2 *Player) *Game {
	g := Game{
		Player1:    p1,
		Player2:    p2,
		CurrPlayer: 1,
		Phase:      "PLACING",
	}

	return &g
}

func (g *Game) PlaceShipForPlayer(p *Player, cell string) bool {
	row, col, err := ConvertCell(cell)
	if err != nil {
		return false
	}

	res := p.Board.PlaceShip(row, col)

	if g.Player1.Board.ShipCount == 5 && g.Player2.Board.ShipCount == 5 {
		g.Phase = "PLAYING"
	}

	return res
}

func (g *Game) FireAtOpponent(firingPlayer *Player, cell string) int {
	var opponent *Player

	if firingPlayer == g.Player1 {
		opponent = g.Player2
	} else {
		opponent = g.Player1
	}

	row, col, err := ConvertCell(cell)

	if err != nil {
		return -1
	}

	res := opponent.Board.Fire(row, col)

	_, gameOver := g.IsGameOver()
	if gameOver {
		g.Phase = "FINISHED"
	}

	return res
}

func (g *Game) SwitchPlayer() int {

	if g.CurrPlayer == 1 {
		g.CurrPlayer = 2
	} else {
		g.CurrPlayer = 1
	}

	return g.CurrPlayer
}

// winner, status
func (g *Game) IsGameOver() (int, bool) {
	if g.Player1.Board.AllShipDestroyed() {
		return 2, true
	}

	if g.Player2.Board.AllShipDestroyed() {
		return 1, true
	}

	return 0, false
}
