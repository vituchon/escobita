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
		mustAdvanceTurn := true
		isComputerTurn := false
		for mustAdvanceTurn {
			if currentRound.HasNextTurn() {
				currentRound.NextTurn()
				isComputerTurn = model.ComputerPlayer.Id == currentRound.CurrentTurnPlayer.Id
				mustAdvanceTurn = !isComputerTurn
			} else {
				if game.CurrentMatch.HasMoreRounds() {
					round := game.CurrentMatch.NextRound()
					round.NextTurn()
					isComputerTurn = model.ComputerPlayer.Id == currentRound.CurrentTurnPlayer.Id
					mustAdvanceTurn = !isComputerTurn
				} else {
					// game ends
					game.CurrentMatch.Ends()
					game.Matchs = append(game.Matchs, *game.CurrentMatch)
					game.CurrentMatch = nil // setting to nil provides a means to detect the current match ending on the client side
					mustAdvanceTurn = false
				}
			}
			if isComputerTurn {
				action := model.CalculateAction(*game.CurrentMatch)
				action, err := game.CurrentMatch.Apply(action)
				if err != nil {
					fmt.Println("ERRORAZO", err)
				} else {
					fmt.Println("Haciendo esta acción", action)
					/* NEED TO MOVE gameswebsockets into service package
					msgPayload := WebSockectOutgoingActionMsgPayload{game, action}
					gameWebSockets.NotifyGameConns(*game.Id, "take", msgPayload)*/
				}

			}
		}
		return &game
	} else {
		// match ended, advancing means to begin another match
		game.CreatesNewMatch() // no need to check the returned error
		// begins the brand new match
		game.CurrentMatch.Prepare()
		game.CurrentMatch.NextRound()             // advances into the first round (within the current match)
		game.CurrentMatch.CurrentRound.NextTurn() // advances into the first turn (within the first round)
		return &game
	}
}
