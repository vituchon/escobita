package model

import (
	"errors"
)

var MatchInProgressErr error = errors.New("A match is in progress")

type Game struct {
	PlayedMatchs []Match  `json:"matchs"`
	Players      []Player `json:"players"`
	CurrentMatch *Match   `json:"currentMatch,omitempty"`
}

func NewGame(players []Player) Game {
	return Game{
		Players:      players,
		PlayedMatchs: make([]Match, 0, 2 /** 36/(len(players)*3) <- TODO: no me acuerdo porque esta formula acá!!**/),
	}
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
	// TODO : Analize if these tree actions could be placed into match within a funtion called "(m *Match) Begins" or "BeginsMatch"
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
