package services

import (
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

func PerformTakeAction(game repositories.PersistentGame, action model.PlayerTakeAction) (*repositories.PersistentGame, model.PlayerAction) {
	updatedAction := game.CurrentMatch.Take(action)
	updatedGame := advanceGame(game)
	return updatedGame, updatedAction
}

func PerformDropAction(game repositories.PersistentGame, action model.PlayerDropAction) (*repositories.PersistentGame, model.PlayerAction) {
	updatedAction := game.CurrentMatch.Drop(action)
	updatedGame := advanceGame(game)
	return updatedGame, updatedAction
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
