package repositories

import (
	"encoding/json"
	"local/escobita/model"
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

type PlayersRepository interface {
	GetPlayers() ([]PersistentPlayer, error)
	GetPlayerById(id int) (*PersistentPlayer, error)
	CreatePlayer(player PersistentPlayer) (created *PersistentPlayer, err error)
	UpdatePlayer(player PersistentPlayer) (updated *PersistentPlayer, err error)
	DeletePlayer(id int) error
}

type PlayersInMemoryRepository struct {
	playersById map[int]PersistentPlayer
}

func NewPlayerInMemoryRepository() *PlayersInMemoryRepository {
	return &PlayersInMemoryRepository{playersById: make(map[int]PersistentPlayer)}
}

func (repo PlayersInMemoryRepository) GetPlayers() ([]PersistentPlayer, error) {
	players := make([]PersistentPlayer, 0, len(repo.playersById))
	for _, player := range repo.playersById {
		players = append(players, player)
	}
	return players, nil
}

func (repo PlayersInMemoryRepository) GetPlayerById(id int) (*PersistentPlayer, error) {
	player, exists := repo.playersById[id]
	if !exists {
		return nil, EntityNotExistsErr
	}
	return &player, nil
}

func (repo *PlayersInMemoryRepository) CreatePlayer(player PersistentPlayer) (created *PersistentPlayer, err error) {
	if player.Id == nil {
		return nil, InvalidEntityStateErr
	}
	repo.playersById[*player.Id] = player
	return &player, nil
}

func (repo *PlayersInMemoryRepository) UpdatePlayer(player PersistentPlayer) (updated *PersistentPlayer, err error) {
	if player.Id == nil {
		return nil, InvalidEntityStateErr
	}
	repo.playersById[*player.Id] = player
	return &player, nil
}

func (repo *PlayersInMemoryRepository) DeletePlayer(id int) error {
	delete(repo.playersById, id)
	return nil
}
