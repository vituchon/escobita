package model

import (
	"errors"

	"github.com/vituchon/escobita/util"
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

func (game *Game) CreatesNewMatch() error {
	if game.CurrentMatch == nil {
		game.createNewMatch()
	} else {
		if game.HasMatchInProgress() {
			return MatchInProgressErr
		}
		game.PlayedMatchs = append(game.PlayedMatchs, *game.CurrentMatch)
		game.createNewMatch()
	}
	return nil
}

func (game Game) HasMatchInProgress() bool {
	return game.CurrentMatch != nil
}

func (game *Game) createNewMatch() {
	match := CreateMatch(game.Players)
	game.CurrentMatch = &match
}

// If already joined then nothing happens... perhaps it could return an error just in case...
func (game *Game) Join(player Player) {
	joinedPlayer := util.Find(game.Players, func(gamePlayer Player) bool { return gamePlayer.Id == player.Id })
	playerNotJoined := joinedPlayer == nil
	if playerNotJoined {
		game.Players = append(game.Players, player)
	}
}
