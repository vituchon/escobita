package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

var clientSessions *sessions.CookieStore

func NewSessionStore(key []byte) {
	clientSessions = sessions.NewCookieStore(key)
}

var clientSequenceId int = 0

func GetOrCreateClientSession(request *http.Request) *sessions.Session {
	clientSession, err := clientSessions.Get(request, "client_session")
	if err != nil {
		fmt.Printf("error while retrieving 'client_session' from session store: %+v\n", err)
	}
	if clientSession.IsNew {
		fmt.Printf("creating new session\n")
		// begin not thread safe
		nextId := clientSequenceId + 1
		clientSession.Values["clientId"] = nextId
		clientSequenceId++
		// end not thread safe
	} else {
		fmt.Printf("using existing session, clientId = %v\n", clientSession.Values["clientId"])
	}
	return clientSession
}

func SaveClientSession(request *http.Request, response http.ResponseWriter, clientSession *sessions.Session) error {
	return clientSessions.Save(request, response, clientSession)
}
