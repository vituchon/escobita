package model

import (
	"errors"

	"github.com/vituchon/escobita/util"
)

var MatchInProgressErr error = errors.New("The match is in progress")
var GameStartedErr error = errors.New("The game is started")
var PlayerAlreadyJoinedErr error = errors.New("The player has already joined the game")
var PlayerNotJoinedErr error = errors.New("The player has not joined the game")

type Game struct {
	PlayedMatchs []Match  `json:"matchs"`
	Players      []Player `json:"players"`
	CurrentMatch *Match   `json:"currentMatch,omitempty"`
}

func NewGame(players []Player) Game {
	return Game{
		Players:      players,
		PlayedMatchs: make([]Match, 0, 2),
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

func (game Game) IsStarted() bool {
	return len(game.PlayedMatchs) > 0 || game.HasMatchInProgress()
}

func (game *Game) createNewMatch() {
	match := CreateMatch(game.Players)
	game.CurrentMatch = &match
}

func (game *Game) Join(player Player) error {
	if game.IsStarted() {
		return GameStartedErr
	}
	joinedPlayer := util.Find(game.Players, func(gamePlayer Player) bool { return gamePlayer.Id == player.Id })
	playerNotJoined := joinedPlayer == nil
	if playerNotJoined {
		game.Players = append(game.Players, player)
	} else {
		return PlayerAlreadyJoinedErr
	}
	return nil
}

func (game *Game) Quit(player Player) error {
	if game.IsStarted() {
		return GameStartedErr // if the game is already started you can quit! as The Eagles stated "you can checkout any time you want, but you can never leave"
	}
	var playerIndex int = -1
	for i, gamePlayer := range game.Players {
		if gamePlayer.Id == player.Id {
			playerIndex = i
			break
		}
	}
	playerJoined := playerIndex != -1
	if playerJoined {
		game.Players = append(game.Players[:playerIndex], game.Players[playerIndex+1:]...) // taken advice from https://github.com/golang/go/wiki/SliceTricks#delete
	} else {
		return PlayerNotJoinedErr
	}
	return nil
}
