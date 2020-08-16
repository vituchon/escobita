package services

// trying to use one file for all kinds services

import (
	"errors"
	"local/escobita/model"
)

var EntityNotExistsErr error = errors.New("Entity doesn't exists")
var EntityDuplicatedErr error = errors.New("Duplicated Entity")

// games

type WebGame struct {
	model.Game `json:"game"`
	Id         *int   `json:"id,omitempty"`
	Name       string `json:"name"`
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
		return nil, EntityDuplicatedErr
	}

	// not treat safe
	nextId := idSequence + 1
	gamesById[idSequence] = game
	game.Id = &nextId
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
