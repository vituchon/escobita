package services

import (
	"encoding/json"
	"github.com/vituchon/escobita/model"
	"strconv"
)

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
