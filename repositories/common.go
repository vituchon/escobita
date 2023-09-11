package repositories

import (
	"errors"
)

var EntityNotExistsErr error = errors.New("Entity doesn't exists")
var DuplicatedEntityErr error = errors.New("Duplicated Entity")
var InvalidEntityStateErr error = errors.New("Entity state is invalid")

type Games interface { // NOTAR queno llamo GamesRepository pues el nombre del paquete sirve de prefijo y al usar estar interfaz desde otro paquete queda `repositories.Games`, ver  @presentation\web\controllers\games.go#gamesRepository
	GetGames() ([]PersistentGame, error)
	GetGameById(id int) (*PersistentGame, error)
	CreateGame(game PersistentGame) (created *PersistentGame, err error)
	UpdateGame(game PersistentGame) (updated *PersistentGame, err error)
	DeleteGame(id int) error
	GetGamesCreatedCount(playerId int) int
}

type Players interface { // aplica misma idea que con Games por eso no agrego el Repository como posfijo
	GetPlayers() ([]PersistentPlayer, error)
	GetPlayerById(id int) (*PersistentPlayer, error)
	CreatePlayer(player PersistentPlayer) (created *PersistentPlayer, err error)
	UpdatePlayer(player PersistentPlayer) (updated *PersistentPlayer, err error)
	DeletePlayer(id int) error
}

type Messages interface {
	GetMessages() ([]PersistentMessage, error)
	GetMessagesByGame(gameId int) ([]PersistentMessage, error)
	GetMessagesByGameAndTime(gameId int, since int64) ([]PersistentMessage, error)
	GetMessageById(id int) (*PersistentMessage, error)
	CreateMessage(message PersistentMessage) (created *PersistentMessage, err error)
	UpdateMessage(message PersistentMessage) (updated *PersistentMessage, err error)
	DeleteMessage(id int) error
}
