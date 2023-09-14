package controllers

import (
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
)

var clientSessions *sessions.CookieStore

func NewSessionStore(key []byte) {
	clientSessions = sessions.NewCookieStore(key)
}

var clientSequenceId int = 0
var mutex sync.Mutex

func GetOrCreateClientSession(request *http.Request) (*sessions.Session, error) {
	clientSession, err := clientSessions.Get(request, "client_session")
	if err != nil {
		return nil, err
	}
	if clientSession.IsNew {
		mutex.Lock()
		defer mutex.Unlock()
		nextId := clientSequenceId + 1
		clientSession.Values["clientId"] = nextId
		clientSequenceId++
	}
	return clientSession, nil
}

func SaveClientSession(request *http.Request, response http.ResponseWriter, clientSession *sessions.Session) error {
	return clientSessions.Save(request, response, clientSession)
}
