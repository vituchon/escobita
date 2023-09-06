package repositories

import (
	"sync"

	"github.com/vituchon/escobita/model"
)

type PersistentGame struct {
	model.Game               // not using json notation intenttonaly in order to marshall the model.Game fields without wrapping into a new subfield
	Id         *int          `json:"id,omitempty"`
	Name       string        `json:"name"`
	PlayerId   int           `json:"playerId"`         // owner
	Matchs     []model.Match `json:"matchs,omitempty"` // played matchs
}

type GamesMemoryStorage struct {
	gamesById  map[int]PersistentGame
	idSequence int
	mutex      sync.Mutex
}

func NewGamesMemoryStorage() *GamesMemoryStorage {
	return &GamesMemoryStorage{gamesById: make(map[int]PersistentGame), idSequence: 0}
}

func (repo GamesMemoryStorage) GetGames() ([]PersistentGame, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	games := make([]PersistentGame, 0, len(repo.gamesById))
	for _, game := range repo.gamesById {
		games = append(games, game)
	}
	return games, nil
}

func (repo GamesMemoryStorage) GetGameById(id int) (*PersistentGame, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	game, exists := repo.gamesById[id]
	if !exists {
		return nil, EntityNotExistsErr
	}
	return &game, nil
}

func (repo *GamesMemoryStorage) CreateGame(game PersistentGame) (created *PersistentGame, err error) {
	if game.Id != nil {
		return nil, DuplicatedEntityErr
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	nextId := repo.idSequence + 1
	game.Id = &nextId
	repo.gamesById[nextId] = game
	repo.idSequence++ // can not reference idSequence as each update would increment all the games Id by id (thus all will be the same)
	return &game, nil
}

func (repo *GamesMemoryStorage) UpdateGame(game PersistentGame) (updated *PersistentGame, err error) {
	if game.Id == nil {
		return nil, EntityNotExistsErr
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	repo.gamesById[*game.Id] = game
	return &game, nil
}

func (repo *GamesMemoryStorage) DeleteGame(id int) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	delete(repo.gamesById, id)
	return nil
}
