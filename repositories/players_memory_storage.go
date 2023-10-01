package repositories

import (
	"sync"

	"github.com/vituchon/escobita/model"
)

type PersistentPlayer = model.Player

type PlayersMemoryRepository struct {
	playersById map[int]PersistentPlayer
	mutex       sync.Mutex
}

func NewPlayersMemoryRepository() *PlayersMemoryRepository {
	return &PlayersMemoryRepository{playersById: make(map[int]PersistentPlayer)}
}

func (repo *PlayersMemoryRepository) GetPlayers() ([]PersistentPlayer, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	players := make([]PersistentPlayer, 0, len(repo.playersById))
	for _, player := range repo.playersById {
		players = append(players, player)
	}
	return players, nil
}

func (repo *PlayersMemoryRepository) GetPlayerById(id int) (*PersistentPlayer, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	player, exists := repo.playersById[id]
	if !exists {
		return nil, EntityNotExistsErr
	}
	return &player, nil
}

func (repo *PlayersMemoryRepository) CreatePlayer(player PersistentPlayer) (created *PersistentPlayer, err error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	repo.playersById[player.Id] = player
	return &player, nil
}

func (repo *PlayersMemoryRepository) UpdatePlayer(player PersistentPlayer) (updated *PersistentPlayer, err error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	repo.playersById[player.Id] = player
	return &player, nil
}

func (repo *PlayersMemoryRepository) DeletePlayer(id int) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	delete(repo.playersById, id)
	return nil
}
