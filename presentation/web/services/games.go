package services

import (
	"errors"
	"fmt"

	"github.com/vituchon/escobita/model"
	"github.com/vituchon/escobita/repositories"
)

// Escobita Oriented Functions

func StartGame(game repositories.PersistentGame) (*repositories.PersistentGame, error) {
	if game.HasMatchInProgress() {
		return nil, model.MatchInProgressErr
	}
	updatedGame := advanceGameComputerAware(game)
	return updatedGame, nil
}

func PerformTakeAction(game repositories.PersistentGame, action model.PlayerTakeAction) (*repositories.PersistentGame, *model.PlayerAction, error) {
	// TODO : move this validation into model facade
	if game.CurrentMatch == nil {
		errMsg := fmt.Sprintf("Can not perform take action: not current match in game(id='%d')", game.Id)
		return nil, nil, errors.New(errMsg)
	}
	// TODO: end logic validation
	updatedAction, err := game.CurrentMatch.Take(action) // TODO: analize if is required to return a copy as the game already is modified due to mutator method invokation
	if err != nil {
		return nil, nil, err
	}
	updatedGame := advanceGameComputerAware(game)
	return updatedGame, &updatedAction, nil
}

func PerformDropAction(game repositories.PersistentGame, action model.PlayerDropAction) (*repositories.PersistentGame, *model.PlayerAction, error) {
	// TODO : move this validation into model facade
	if game.CurrentMatch == nil {
		errMsg := fmt.Sprintf("Can not perform drop action: not current match in game(id='%d')", game.Id)
		return nil, nil, errors.New(errMsg)
	}
	// TODO: end logic validation
	updatedAction, err := game.CurrentMatch.Drop(action)
	if err != nil {
		return nil, nil, err
	}
	updatedGame := advanceGameComputerAware(game)
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
		game.CreatesNewMatch() // no need to check the returned error
		// begins the brand new match
		game.CurrentMatch.Prepare()
		game.CurrentMatch.NextRound()             // advances into the first round (within the current match)
		game.CurrentMatch.CurrentRound.NextTurn() // advances into the first turn (within the first round)
		return &game
	}
}

func advanceGameComputerAware(game repositories.PersistentGame) *repositories.PersistentGame {
	var updated *repositories.PersistentGame = &game
	mustResume := true
	for mustResume {
		updated = advanceGame(*updated)
		if updated.CurrentMatch != nil { // TODO: method for determinging is the current's game match is not ended.. if any remaining player must act
			isComputerTurn := model.ComputerPlayer.Id == updated.CurrentMatch.CurrentRound.CurrentTurnPlayer.Id
			fmt.Println("updated.CurrentMatch.CurrentRound.CurrentTurnPlayer", updated.CurrentMatch.CurrentRound.CurrentTurnPlayer)
			fmt.Println("isComputerTurn", isComputerTurn)
			mustResume = isComputerTurn //  must resume UNTIL match ends or there is an human player turn
			if isComputerTurn {
				action := model.CalculateAction(*updated.CurrentMatch)
				action, _ = updated.CurrentMatch.Apply(action)
				msgPayload := WebSockectOutgoingActionMsgPayload{updated, &action}
				switch action.(type) {
				case model.PlayerTakeAction:
					GameWebSockets.NotifyGameConns(*game.Id, "take", msgPayload)
					break
				case model.PlayerDropAction:
					GameWebSockets.NotifyGameConns(*game.Id, "drop", msgPayload)
					break
				}
			}
		} else {
			mustResume = false
		}
	}
	return updated

}
