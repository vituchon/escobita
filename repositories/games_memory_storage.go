package repositories

import (
	"sync"

	"github.com/vituchon/escobita/model"
)

type PersistentGame struct {
	model.Game                  // not using json notation intenttonaly in order to marshall the model.Game fields without wrapping into a new subfield
	Id         *int             `json:"id,omitempty"`
	Name       string           `json:"name"`
	Owner      PersistentPlayer `json:"owner"`
	Matchs     []model.Match    `json:"matchs,omitempty"` // played matchs
}

type GamesMemoryRepository struct {
	gamesById              map[int]PersistentGame
	gamesCreatedByPlayerId map[int]int
	idSequence             int
	mutex                  sync.Mutex
}

func NewGamesMemoryRepository() *GamesMemoryRepository {
	return &GamesMemoryRepository{gamesById: make(map[int]PersistentGame), gamesCreatedByPlayerId: make(map[int]int), idSequence: 0}
}

func (repo *GamesMemoryRepository) GetGames() ([]PersistentGame, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	games := make([]PersistentGame, 0, len(repo.gamesById))
	for _, game := range repo.gamesById {
		games = append(games, game)
	}
	return games, nil
}

func (repo *GamesMemoryRepository) GetGameById(id int) (*PersistentGame, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	game, exists := repo.gamesById[id]
	if !exists {
		return nil, EntityNotExistsErr
	}
	return &game, nil
}

func (repo *GamesMemoryRepository) CreateGame(game PersistentGame) (created *PersistentGame, err error) {
	if game.Id != nil {
		return nil, DuplicatedEntityErr
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	nextId := repo.idSequence + 1 // need to copy repo.idSequence in another place (nextId), also added plus one to increment the sequence current number...
	game.Id = &nextId             // ...that place (nextId) will work as reference...
	repo.gamesById[nextId] = game
	repo.idSequence++ // ...if it is used idSequence as a reference, then each update would increment all the games Id by 1 (actually all game.Id will point to the same thing)
	repo.gamesCreatedByPlayerId[game.Owner.Id]++
	return &game, nil
}

func (repo *GamesMemoryRepository) UpdateGame(game PersistentGame) (updated *PersistentGame, err error) {
	if game.Id == nil {
		return nil, EntityNotExistsErr
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	repo.gamesById[*game.Id] = game
	return &game, nil
}

func (repo *GamesMemoryRepository) DeleteGame(id int) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	game := repo.gamesById[id]
	repo.gamesCreatedByPlayerId[game.Owner.Id]--
	delete(repo.gamesById, id)
	return nil
}

func (repo GamesMemoryRepository) GetGamesCreatedCount(playerId int) int {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	return repo.gamesCreatedByPlayerId[playerId]
}
