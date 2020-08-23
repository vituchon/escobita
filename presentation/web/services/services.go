package services

// trying to use one file for all kinds services

// Contains code designed to be used in gui web-client interface.

import (
	"encoding/json"
	"errors"
	"local/escobita/model"
	"strconv"
)

var EntityNotExistsErr error = errors.New("Entity doesn't exists")
var DuplicatedEntityErr error = errors.New("Duplicated Entity")
var InvalidEntityStateErr error = errors.New("Entity state is invalid")

// GAMES

type WebGame struct {
	model.Game        // not using json notation intenttonaly in order to marshall the model.Game fields without wrapping into a new subfield
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

// advances the game into his next state, that is, a new match or a new round or ends
func AdvanceGame(id int) (*WebGame, error) {
	game, exists := gamesById[id]
	if !exists {
		return nil, EntityNotExistsErr
	}
	if game.HasMatchInProgress() {
		currentRound := game.CurrentMatch.CurrentRound
		if currentRound.HasNextTurn() {
			currentRound.NextTurn()
		} else {
			if game.CurrentMatch.HasMoreRounds() {
				game.CurrentMatch.NextRound()
			} else {
				// game ends
				game.CurrentMatch.Ends()
			}
		}
		return &game, nil
	} else {
		err := game.BeginMatch()
		return &game, err
	}

}

// PLAYERS
type WebPlayer struct {
	model.Player
	Id *int `json:"id"`
}

func (wp WebPlayer) String() string {
	if wp.Id == nil {
		return "NO_ID " + wp.Name
	}
	return strconv.Itoa(*wp.Id) + " " + wp.Name
}

func (wp WebPlayer) MarshalJSON() ([]byte, error) {
	if wp.Id == nil {
		return []byte(`{"name":"` + wp.Name + `"}`), nil
	}
	return []byte(`{"name":"` + wp.Name + `", "id":` + strconv.Itoa(*wp.Id) + `}`), nil
}

func (wp *WebPlayer) UnmarshalJSON(b []byte) error {
	var stuff map[string]interface{}
	err := json.Unmarshal(b, &stuff)
	if err != nil {
		return err
	}
	wp.Name = stuff["name"].(string)
	id := int(stuff["id"].(float64))
	wp.Id = &id
	return nil
}

var playersById map[int]WebPlayer = make(map[int]WebPlayer)

func GetPlayers() ([]WebPlayer, error) {
	players := make([]WebPlayer, 0, len(playersById))
	for _, player := range playersById {
		players = append(players, player)
	}
	return players, nil
}

func GetPlayerById(id int) (*WebPlayer, error) {
	player, exists := playersById[id]
	if !exists {
		return nil, EntityNotExistsErr
	}
	return &player, nil
}

func CreatePlayer(player WebPlayer) (created *WebPlayer, err error) {
	if player.Id == nil {
		return nil, InvalidEntityStateErr
	}
	playersById[*player.Id] = player
	return &player, nil
}

func UpdatePlayer(player WebPlayer) (updated *WebPlayer, err error) {
	if player.Id == nil {
		return nil, InvalidEntityStateErr
	}
	playersById[*player.Id] = player
	return &player, nil
}

func DeletePlayer(id int) error {
	delete(playersById, id)
	return nil
}
