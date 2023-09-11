package services

import (
	"github.com/vituchon/escobita/repositories"
)

type WebMessage = repositories.PersistentMessage

var messagesRepository repositories.Messages

func init() {
	messagesRepository = repositories.NewMessagesMemoryStorage()
}

func GetMessages() ([]WebMessage, error) {
	return messagesRepository.GetMessages()
}

func GetMessagesByGame(gameId int) ([]WebMessage, error) {
	return messagesRepository.GetMessagesByGame(gameId)
}

func GetMessagesByGameAndTime(gameId int, since int64) ([]WebMessage, error) {
	return messagesRepository.GetMessagesByGameAndTime(gameId, since)
}

func GetMessageById(id int) (*WebMessage, error) {
	return messagesRepository.GetMessageById(id)
}

func CreateMessage(message WebMessage) (created *WebMessage, err error) {
	return messagesRepository.CreateMessage(message)
}

func UpdateMessage(message WebMessage) (updated *WebMessage, err error) {
	return messagesRepository.UpdateMessage(message)
}

func DeleteMessage(id int) error {
	return messagesRepository.DeleteMessage(id)
}
