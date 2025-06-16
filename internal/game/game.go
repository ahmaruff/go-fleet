package game

type Game struct {
	Player1    *Player
	Player2    *Player
	CurrPlayer int
	Phase      string
}
