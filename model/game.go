package model

type Game struct {
	Matchs         []Match        `json:"matchs"`
	Players        []Player       `json:"players"`
	ScorePerPlayer map[Player]int `json:"scorePerPlayer"`
}

func NewGame(players []Player) Game {
	game := Game{
		Players:        players,
		ScorePerPlayer: make(map[Player]int),
		Matchs:         make([]Match, 0, 36/(len(players)*3)),
	}
	for _, player := range players {
		game.ScorePerPlayer[player] = 0
	}
	return game
}
