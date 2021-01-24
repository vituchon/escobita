package services

import (
	"sync"
	"time"
)

type WebMessage struct {
	Id       *int   `json:"id,omitempty"`
	PlayerId int    `json:"playerId"`
	GameId   int    `json:"gameId"`
	Text     string `json:"text"`
	Created  int64  `json:"created"`
}

var messagesById map[int]WebMessage = make(map[int]WebMessage)

func GetMessages() ([]WebMessage, error) {
	messages := make([]WebMessage, 0, len(messagesById))
	for _, message := range messagesById {
		messages = append(messages, message)
	}
	return messages, nil
}

// the game would serve as the room
// dev notes: GetMessages y GetMessagesByGame son funciones iguales con filtro disinto, getMessages es filtro TRUE y el
// otro por match de GameId, por ahora son busquedas sequenciales
func GetMessagesByGame(gameId int) ([]WebMessage, error) {
	messages := make([]WebMessage, 0, len(messagesById))
	for _, message := range messagesById {
		if message.GameId == gameId {
			messages = append(messages, message)
		}
	}
	return messages, nil
}

func GetMessageById(id int) (*WebMessage, error) {
	message, exists := messagesById[id]
	if !exists {
		return nil, EntityNotExistsErr
	}
	return &message, nil
}

var idMessageSequence int = 0
var mutex = &sync.Mutex{}

func CreateMessage(message WebMessage) (created *WebMessage, err error) {
	if message.Id != nil {
		return nil, InvalidEntityStateErr
	}

	mutex.Lock()
	nextId := idMessageSequence + 1
	message.Id = &nextId
	messagesById[nextId] = message
	idMessageSequence++ // can not reference idSequence as each update would increment all the games Id by id (thus all will be the same)
	mutex.Unlock()

	message.Created = time.Now().Unix()
	messagesById[*message.Id] = message
	return &message, nil
}

func UpdateMessage(message WebMessage) (updated *WebMessage, err error) {
	if message.Id == nil {
		return nil, InvalidEntityStateErr
	}
	messagesById[*message.Id] = message
	return &message, nil
}

func DeleteMessage(id int) error {
	delete(messagesById, id)
	return nil
}
