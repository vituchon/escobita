package repositories

import (
	"encoding/json"
	"github.com/vituchon/escobita/model"
	"strconv"
)

type PersistentPlayer struct {
	model.Player
	Id *int `json:"id"`
}

func (pp PersistentPlayer) String() string {
	if pp.Id == nil {
		return "NO_ID " + pp.Name
	}
	return strconv.Itoa(*pp.Id) + " " + pp.Name
}

func (pp PersistentPlayer) MarshalJSON() ([]byte, error) {
	if pp.Id == nil {
		return []byte(`{"name":"` + pp.Name + `"}`), nil
	}
	return []byte(`{"name":"` + pp.Name + `", "id":` + strconv.Itoa(*pp.Id) + `}`), nil
}

func (pp *PersistentPlayer) UnmarshalJSON(b []byte) error {
	var stuff map[string]interface{}
	err := json.Unmarshal(b, &stuff)
	if err != nil {
		return err
	}
	pp.Name = stuff["name"].(string)
	id := int(stuff["id"].(float64))
	pp.Id = &id
	return nil
}

type Players interface {
	GetPlayers() ([]PersistentPlayer, error)
	GetPlayerById(id int) (*PersistentPlayer, error)
	CreatePlayer(player PersistentPlayer) (created *PersistentPlayer, err error)
	UpdatePlayer(player PersistentPlayer) (updated *PersistentPlayer, err error)
	DeletePlayer(id int) error
}

type PlayersMemoryStorage struct {
	playersById map[int]PersistentPlayer
}

func NewPlayersMemoryStorage() *PlayersMemoryStorage {
	return &PlayersMemoryStorage{playersById: make(map[int]PersistentPlayer)}
}

func (repo PlayersMemoryStorage) GetPlayers() ([]PersistentPlayer, error) {
	players := make([]PersistentPlayer, 0, len(repo.playersById))
	for _, player := range repo.playersById {
		players = append(players, player)
	}
	return players, nil
}

func (repo PlayersMemoryStorage) GetPlayerById(id int) (*PersistentPlayer, error) {
	player, exists := repo.playersById[id]
	if !exists {
		return nil, EntityNotExistsErr
	}
	return &player, nil
}

func (repo *PlayersMemoryStorage) CreatePlayer(player PersistentPlayer) (created *PersistentPlayer, err error) {
	if player.Id == nil {
		return nil, InvalidEntityStateErr
	}
	repo.playersById[*player.Id] = player
	return &player, nil
}

func (repo *PlayersMemoryStorage) UpdatePlayer(player PersistentPlayer) (updated *PersistentPlayer, err error) {
	if player.Id == nil {
		return nil, InvalidEntityStateErr
	}
	repo.playersById[*player.Id] = player
	return &player, nil
}

func (repo *PlayersMemoryStorage) DeletePlayer(id int) error {
	delete(repo.playersById, id)
	return nil
}
