package repositories

import (
	"sync"
	"time"
)

type PersistentMessage struct {
	Id       *int   `json:"id,omitempty"`
	PlayerId int    `json:"playerId"`
	GameId   int    `json:"gameId"`
	Text     string `json:"text"`
	Created  int64  `json:"created"`
}

type MessagesMemoryStorage struct {
	messagesById map[int]PersistentMessage
	mutex        sync.Mutex
	idSequence   int
}

func NewMessagesMemoryStorage() *MessagesMemoryStorage {
	return &MessagesMemoryStorage{messagesById: make(map[int]PersistentMessage)}
}

func (repo *MessagesMemoryStorage) GetMessages() ([]PersistentMessage, error) {
	return repo.doGetMessages(MessageFilterByNone)
}

func (repo *MessagesMemoryStorage) GetMessagesByGame(gameId int) ([]PersistentMessage, error) {
	filter := MessageFilterByGameId{gameId}
	return repo.doGetMessages(filter.fulfill)
}

func (repo *MessagesMemoryStorage) GetMessagesByGameAndTime(gameId int, since int64) ([]PersistentMessage, error) {
	filterByGame := MessageFilterByGameId{gameId}
	filterByTime := MessageFilterByTime{since}
	filterByBoth := func(message PersistentMessage) bool {
		return filterByGame.fulfill(message) && filterByTime.fulfill(message)
	}
	return repo.doGetMessages(filterByBoth)
}

func (repo *MessagesMemoryStorage) doGetMessages(filterFunc MessageFilterFunc) ([]PersistentMessage, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	messages := make([]PersistentMessage, 0, len(repo.messagesById))
	for _, message := range repo.messagesById {
		if filterFunc(message) {
			messages = append(messages, message)
		}
	}
	return messages, nil
}

func (repo *MessagesMemoryStorage) GetMessageById(id int) (*PersistentMessage, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	message, exists := repo.messagesById[id]
	if !exists {
		return nil, EntityNotExistsErr
	}
	return &message, nil
}

func (repo *MessagesMemoryStorage) CreateMessage(message PersistentMessage) (created *PersistentMessage, err error) {
	if message.Id != nil {
		return nil, InvalidEntityStateErr
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	nextId := repo.idSequence + 1
	message.Id = &nextId
	repo.messagesById[nextId] = message
	repo.idSequence++ // can not reference idSequence as each update would increment all the games Id by id (thus all will be the same)

	message.Created = time.Now().Unix()
	repo.messagesById[*message.Id] = message
	return &message, nil
}

func (repo *MessagesMemoryStorage) UpdateMessage(message PersistentMessage) (updated *PersistentMessage, err error) {
	if message.Id == nil {
		return nil, InvalidEntityStateErr
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	repo.messagesById[*message.Id] = message
	return &message, nil
}

func (repo *MessagesMemoryStorage) DeleteMessage(id int) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	delete(repo.messagesById, id)
	return nil
}

type MessageFilterFunc = func(message PersistentMessage) bool

func MessageFilterByNone(message PersistentMessage) bool {
	return true
}

type MessageFilterByGameId struct {
	GameId int
}

func (f MessageFilterByGameId) fulfill(message PersistentMessage) bool {
	return message.GameId == f.GameId
}

type MessageFilterByTime struct {
	Since int64
}

func (f MessageFilterByTime) fulfill(message PersistentMessage) bool {
	return message.Created >= f.Since
}
