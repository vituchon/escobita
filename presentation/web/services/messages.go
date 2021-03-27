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

type MessageFilterFunc = func(message WebMessage) bool

func MessageFilterByNone(message WebMessage) bool {
	return true
}

type MessageFilterByGameId struct {
	GameId int
}

func (f MessageFilterByGameId) fulfill(message WebMessage) bool {
	return message.GameId == f.GameId
}

type MessageFilterByTime struct {
	Since int64
}

func (f MessageFilterByTime) fulfill(message WebMessage) bool {
	return message.Created >= f.Since
}

func doGetMessages(filterFunc MessageFilterFunc) ([]WebMessage, error) {
	messages := make([]WebMessage, 0, len(messagesById))
	for _, message := range messagesById {
		if filterFunc(message) {
			messages = append(messages, message)
		}
	}
	return messages, nil
}

func GetMessages() ([]WebMessage, error) {
	return doGetMessages(MessageFilterByNone)
}

// the game would serve as the "room"
func GetMessagesByGame(gameId int) ([]WebMessage, error) {
	filter := MessageFilterByGameId{gameId}
	return doGetMessages(filter.fulfill)
}

func GetMessagesByGameAndTime(gameId int, since int64) ([]WebMessage, error) {
	filterByGame := MessageFilterByGameId{gameId}
	filterByTime := MessageFilterByTime{since}
	filterByBoth := func(message WebMessage) bool {
		return filterByGame.fulfill(message) && filterByTime.fulfill(message)
	}
	return doGetMessages(filterByBoth)
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
