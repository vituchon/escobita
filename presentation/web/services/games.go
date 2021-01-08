package services

import (
	"local/escobita/model"
)

type WebGame struct {
	model.Game               // not using json notation intenttonaly in order to marshall the model.Game fields without wrapping into a new subfield
	Id         *int          `json:"id,omitempty"`
	Name       string        `json:"name"`
	PlayerId   int           `json:"playerId"`          // owner
	Matchs     []model.Match `json:"matchs, omitempty"` // played matchs
}

// meet the actual storate...
var gamesById map[int]WebGame = make(map[int]WebGame)
var idSequence int = 0

// and his basic interface..
func GetGames() ([]WebGame, error) {
	games := make([]WebGame, 0, len(gamesById))
	for _, game := range gamesById {
		games = append(games, game)
	}
	return games, nil
}

func GetGameById(id int) (*WebGame, error) {
	game, exists := gamesById[id]
	if !exists {
		return nil, EntityNotExistsErr
	}
	return &game, nil
}

func CreateGame(game WebGame) (created *WebGame, err error) {
	if game.Id != nil {
		return nil, DuplicatedEntityErr
	}
	// not treat safe
	nextId := idSequence + 1
	game.Id = &nextId
	gamesById[nextId] = game
	idSequence++ // can not reference idSequence as each update would increment all the games Id by id (thus all will be the same)
	// end not treat safe
	return &game, nil
}

func UpdateGame(game WebGame) (updated *WebGame, err error) {
	if game.Id == nil {
		return nil, EntityNotExistsErr
	}
	gamesById[*game.Id] = game
	return &game, nil
}

func DeleteGame(id int) error {
	delete(gamesById, id)
	return nil
}

// Escobita Oriented Functions

func ResumeGame(game WebGame) (*WebGame, error) {
	if game.HasMatchInProgress() {
		return nil, model.MatchInProgressErr
	}
	updatedGame, err := advanceGame(game)
	return updatedGame, err
}

func PerformTakeAction(game WebGame, action model.PlayerTakeAction) (*WebGame, model.PlayerAction, error) {
	updatedAction := game.CurrentMatch.Take(action)
	updatedGame, err := advanceGame(game)
	return updatedGame, updatedAction, err
}

func PerformDropAction(game WebGame, action model.PlayerDropAction) (*WebGame, model.PlayerAction, error) {
	updatedAction := game.CurrentMatch.Drop(action)
	updatedGame, err := advanceGame(game)
	return updatedGame, updatedAction, err
}

func CalculateCurrentMatchStats(game WebGame) model.ScoreSummaryByPlayer {
	staticticsByPlayer := game.CurrentMatch.CalculateStaticticsByPlayer()
	scoreSummaryByPlayer := staticticsByPlayer.BuildScoreSummaryByPlayer()
	return scoreSummaryByPlayer
}

func CalculatePlayedMatchStats(game WebGame, index int) model.ScoreSummaryByPlayer {
	staticticsByPlayer := game.Matchs[index].CalculateStaticticsByPlayer()
	scoreSummaryByPlayer := staticticsByPlayer.BuildScoreSummaryByPlayer()
	return scoreSummaryByPlayer
}

// advances the game into his next state, that is, a new match or a new round or ends
func advanceGame(game WebGame) (*WebGame, error) {
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
		_, err := UpdateGame(game)
		return &game, err
	} else {
		// match ended, advancing means to begin another match
		game.BeginsNewMatch() // no need to check the returned error
		_, err := UpdateGame(game)
		return &game, err

	}
}
