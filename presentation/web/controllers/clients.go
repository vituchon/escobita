package controllers

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/vituchon/escobita/util"
)

var clientSessions *sessions.CookieStore
var integerSequence util.IntegerSequence

func InitSessionStore(key []byte) {
	clientSessions = sessions.NewCookieStore(key)
	integerSequence = util.NewFsIntegerSequence("escobita.seq", 0, 1)
}

func GetOrCreateClientSession(request *http.Request) (*sessions.Session, error) {
	clientSession, err := clientSessions.Get(request, "escoba_client")
	if err != nil {
		return nil, err
	}
	if clientSession.IsNew {
		scaledNextId, err := integerSequence.GetNext()
		if err != nil {
			return nil, err
		}
		clientSession.Values["clientId"] = scaledNextId
	}
	return clientSession, nil
}

func SaveClientSession(request *http.Request, response http.ResponseWriter, clientSession *sessions.Session) error {
	return clientSessions.Save(request, response, clientSession)
}
