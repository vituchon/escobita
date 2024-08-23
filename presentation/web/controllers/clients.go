package controllers

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/vituchon/escobita/presentation/web/services"
)

func GetOrCreateClientSession(request *http.Request, response http.ResponseWriter) (clientSession *services.ClientSession, err error) {
	clientSession, err = getClientSession(request)
	if err != nil {
		isNewClient := (err == http.ErrNoCookie)
		if isNewClient {
			uuid := uuid.New()
			uuidBytes := uuid[:]
			nextId := int(big.NewInt(0).SetBytes(uuidBytes[:8]).Int64())
			clientSession = &services.ClientSession{
				Id: nextId,
			}
			setClientSession(response, clientSession)
			return clientSession, nil
		} else {
			return nil, err
		}
	}
	return clientSession, nil
}

const cookieName string = "escobita_client"

func getClientSession(request *http.Request) (*services.ClientSession, error) {
	cookie, err := request.Cookie(cookieName)
	if err == nil {
		clientIdAsStr := cookie.Value
		clientId, err := strconv.Atoi(clientIdAsStr)
		if err != nil {
			return nil, fmt.Errorf("Client id(=%v) is not an int", clientIdAsStr)
		}
		clientSession := services.ClientSession{
			Id: clientId,
		}
		return &clientSession, nil
	}
	return nil, err
}

func setClientSession(response http.ResponseWriter, clientSession *services.ClientSession) {
	cookieValue := strconv.Itoa(clientSession.Id)

	cookie := http.Cookie{
		Name:  cookieName,
		Value: cookieValue,
	}

	http.SetCookie(response, &cookie)
	return
}
