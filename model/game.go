package model

import (
	"errors"
)

var MatchInProgressErr error = errors.New("A match is in progress")

type Game struct {
	PlayedMatchs   []Match        `json:"matchs"`
	Players        []Player       `json:"players"`
	ScorePerPlayer map[Player]int `json:"scorePerPlayer"` // TODO : rename to ScoreByPlayer
	CurrentMatch   *Match         `json:"CurrentMatch, omitempty"`
}

func NewGame(players []Player) Game {
	game := Game{
		Players:        players,
		ScorePerPlayer: make(map[Player]int),
		PlayedMatchs:   make([]Match, 0, 2 /** 36/(len(players)*3) <- TODO: no me acuerdo porque esta formula acá!!**/),
	}
	for _, player := range players {
		game.ScorePerPlayer[player] = 0
	}
	return game
}

func (game *Game) BeginMatch() error {
	if game.CurrentMatch == nil {
		game.beginNewMatch()
		return nil
	} else {
		if game.HasMatchInProgress() {
			return MatchInProgressErr
		}
		game.PlayedMatchs = append(game.PlayedMatchs, *game.CurrentMatch)
		game.beginNewMatch()
		return nil
	}
}

func (game Game) HasMatchInProgress() bool {
	return game.CurrentMatch.HasMoreRounds()
}

func (game *Game) beginNewMatch() {
	match := CreateMatch(game.Players)
	game.CurrentMatch = &match
}
