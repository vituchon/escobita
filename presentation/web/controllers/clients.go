package controllers

import (
	"net/http"

	"math/big"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
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
		uuid := uuid.New()
		uuidBytes := uuid[:]
		nextId := int(big.NewInt(0).SetBytes(uuidBytes[:8]).Int64())
		min := -1000000000
		max := 1000000000
		scaledNextId := (nextId % (max - min + 1)) + min // los escalo pues los navegadores pueden no soportar números tan grandes y luego hay conflicto a la hora de actualziar el jugador cuando se registra el nombre
		clientSession.Values["clientId"] = scaledNextId
	}
	return clientSession, nil
}

func SaveClientSession(request *http.Request, response http.ResponseWriter, clientSession *sessions.Session) error {
	return clientSessions.Save(request, response, clientSession)
}
