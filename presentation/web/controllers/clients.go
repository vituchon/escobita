package controllers

import (
	"net/http"
	"github.com/gorilla/sessions"
	"github.com/vituchon/escobita/presentation/util"
)

var clientSessions *sessions.CookieStore

func InitSessionStore(key []byte) {
	clientSessions = sessions.NewCookieStore(key)
}

func GetOrCreateClientSession(request *http.Request) (*sessions.Session, error) {
	clientSession, err := clientSessions.Get(request, "escoba_client")
	if err != nil {
		return nil, err
	}
	if clientSession.IsNew {
		scaledNextId := util.GenerateRandomNumber(0, 10000000000)
		clientSession.Values["clientId"] = scaledNextId
	}
	return clientSession, nil
}

func SaveClientSession(request *http.Request, response http.ResponseWriter, clientSession *sessions.Session) error {
	return clientSessions.Save(request, response, clientSession)
}
