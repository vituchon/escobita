package model

import (
	"errors"
)

var MatchInProgressErr error = errors.New("A match is in progress")

type Game struct {
	PlayedMatchs  []Match        `json:"matchs"`
	Players       []Player       `json:"players"`
	ScoreByPlayer map[Player]int `json:"scoreByPlayerName"` // dev notes (TODO): The map uses the id+name... not only name, shall I change the name?
	CurrentMatch  *Match         `json:"currentMatch,omitempty"`
}

func NewGame(players []Player) Game {
	game := Game{
		Players:       players,
		ScoreByPlayer: make(map[Player]int),
		PlayedMatchs:  make([]Match, 0, 2 /** 36/(len(players)*3) <- TODO: no me acuerdo porque esta formula acá!!**/),
	}
	for _, player := range players {
		game.ScoreByPlayer[player] = 0
	}
	return game
}

func (game *Game) BeginsNewMatch() error {
	if game.CurrentMatch == nil {
		game.createNewMatch()
	} else {
		if game.HasMatchInProgress() {
			return MatchInProgressErr
		}
		game.PlayedMatchs = append(game.PlayedMatchs, *game.CurrentMatch)
		game.createNewMatch()
	}
	// there is a current match after executing above statements
	game.CurrentMatch.Begins()
	game.CurrentMatch.NextRound()
	// there is a current round after executing above statements
	game.CurrentMatch.CurrentRound.NextTurn()
	return nil
}

func (game Game) HasMatchInProgress() bool {
	return game.CurrentMatch != nil
}

func (game *Game) createNewMatch() {
	match := CreateMatch(game.Players)
	game.CurrentMatch = &match
}
