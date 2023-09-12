package services

import (
	"errors"

	"github.com/vituchon/escobita/repositories"
)

type VolatileWebMessage struct {
	GameId int                           `json:"gameId"`
	Player repositories.PersistentPlayer `json:"player"`
	Text   string                        `json:"text"`
}

type PersistentWebMessage = repositories.PersistentMessage

type VolatileMessageRepository struct {
}

var NotPersistentMessageRepositoryErr = errors.New("Can not perform the given operation with a the volatile message repository. Try using MessagesMemoryRepository.")

func (repo *VolatileMessageRepository) GetMessages() ([]repositories.PersistentMessage, error) {
	return nil, NotPersistentMessageRepositoryErr
}

func (repo *VolatileMessageRepository) GetMessagesByGame(gameId int) ([]repositories.PersistentMessage, error) {
	return nil, NotPersistentMessageRepositoryErr
}

func (repo *VolatileMessageRepository) GetMessagesByGameAndTime(gameId int, since int64) ([]repositories.PersistentMessage, error) {
	return nil, NotPersistentMessageRepositoryErr
}

func (repo *VolatileMessageRepository) GetMessageById(id int) (*repositories.PersistentMessage, error) {
	return nil, NotPersistentMessageRepositoryErr
}

func (repo *VolatileMessageRepository) CreateMessage(message repositories.PersistentMessage) (created *repositories.PersistentMessage, err error) {
	return nil, NotPersistentMessageRepositoryErr
}

func (repo *VolatileMessageRepository) UpdateMessage(message repositories.PersistentMessage) (updated *repositories.PersistentMessage, err error) {
	return nil, NotPersistentMessageRepositoryErr
}

func (repo *VolatileMessageRepository) DeleteMessage(id int) error {
	return NotPersistentMessageRepositoryErr
}

var messagesRepository repositories.Messages

func init() {
	//messagesRepository = repositories.NewMessagesMemoryRepository()
	messagesRepository = &VolatileMessageRepository{}
}

func GetMessages() ([]repositories.PersistentMessage, error) {
	return messagesRepository.GetMessages()
}

func GetMessagesByGame(gameId int) ([]repositories.PersistentMessage, error) {
	return messagesRepository.GetMessagesByGame(gameId)
}

func GetMessagesByGameAndTime(gameId int, since int64) ([]repositories.PersistentMessage, error) {
	return messagesRepository.GetMessagesByGameAndTime(gameId, since)
}

func GetMessageById(id int) (*repositories.PersistentMessage, error) {
	return messagesRepository.GetMessageById(id)
}

func CreateMessage(message repositories.PersistentMessage) (created *repositories.PersistentMessage, err error) {
	return messagesRepository.CreateMessage(message)
}

func UpdateMessage(message repositories.PersistentMessage) (updated *repositories.PersistentMessage, err error) {
	return messagesRepository.UpdateMessage(message)
}

func DeleteMessage(id int) error {
	return messagesRepository.DeleteMessage(id)
}
