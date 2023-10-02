package services

import (
	"github.com/vituchon/escobita/model"
)

type WebPlayer = model.Player

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
	playersById[player.Id] = player
	return &player, nil
}

func UpdatePlayer(player WebPlayer) (updated *WebPlayer, err error) {
	playersById[player.Id] = player
	return &player, nil
}

func DeletePlayer(id int) error {
	delete(playersById, id)
	return nil
}
