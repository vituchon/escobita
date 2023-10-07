package services

import (
	"errors"
	"fmt"

	"github.com/vituchon/escobita/model"
	"github.com/vituchon/escobita/repositories"
)

// Escobita Oriented Functions

func ResumeGame(game repositories.PersistentGame) (*repositories.PersistentGame, error) {
	if game.HasMatchInProgress() {
		return nil, model.MatchInProgressErr
	}
	updatedGame := advanceGame(game)
	return updatedGame, nil
}

func PerformTakeAction(game repositories.PersistentGame, action model.PlayerTakeAction) (*repositories.PersistentGame, *model.PlayerAction, error) {
	if game.CurrentMatch == nil {
		errMsg := fmt.Sprintf("Can not perform take action: not current match in game(id=('%d')", game.Id)
		return nil, nil, errors.New(errMsg)
	}
	updatedAction, err := game.CurrentMatch.Take(action)
	if err != nil {
		return nil, nil, err
	}
	updatedGame := advanceGame(game)
	return updatedGame, &updatedAction, nil
}

func PerformDropAction(game repositories.PersistentGame, action model.PlayerDropAction) (*repositories.PersistentGame, *model.PlayerAction, error) {
	if game.CurrentMatch == nil {
		errMsg := fmt.Sprintf("Can not perform drop action: not current match in game(id=('%d')", game.Id)
		return nil, nil, errors.New(errMsg)
	}
	updatedAction, err := game.CurrentMatch.Drop(action)
	if err != nil {
		return nil, nil, err
	}
	updatedGame := advanceGame(game)
	return updatedGame, &updatedAction, nil
}

func CalculateCurrentMatchStats(game repositories.PersistentGame) model.ScoreSummaryByPlayer {
	staticticsByPlayer := game.CurrentMatch.CalculateStaticticsByPlayer()
	scoreSummaryByPlayer := staticticsByPlayer.BuildScoreSummaryByPlayer()
	return scoreSummaryByPlayer
}

func CalculatePlayedMatchStats(game repositories.PersistentGame, index int) model.ScoreSummaryByPlayer {
	staticticsByPlayer := game.Matchs[index].CalculateStaticticsByPlayer()
	scoreSummaryByPlayer := staticticsByPlayer.BuildScoreSummaryByPlayer()
	return scoreSummaryByPlayer
}

func CanPlayerDeleteGame(game *repositories.PersistentGame, player repositories.PersistentPlayer) bool {
	return game.Owner.Id == player.Id
}

// advances the game into his next state, that is, a new match or a new round or ends
func advanceGame(game repositories.PersistentGame) *repositories.PersistentGame {
	if game.HasMatchInProgress() {
		currentRound := game.CurrentMatch.CurrentRound
		if currentRound.HasNextTurn() {
			currentRound.NextTurn()
		} else {
			if game.CurrentMatch.HasMoreRounds() {
				round := game.CurrentMatch.NextRound()
				round.NextTurn()
			} else {
				// game ends
				game.CurrentMatch.Ends()
				game.Matchs = append(game.Matchs, *game.CurrentMatch)
				game.CurrentMatch = nil // setting to nil provides a means to detect the current match ending on the client side
			}
		}
		return &game
	} else {
		// match ended, advancing means to begin another match
		game.BeginsNewMatch() // no need to check the returned error
		return &game
	}
}
