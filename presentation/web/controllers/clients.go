package controllers

import (
	"fmt"
	//"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/vituchon/escobita/presentation/web/services"
)

var clientSequenceId int = 0
var mutex sync.Mutex

func GetOrCreateClientSession(request *http.Request, response http.ResponseWriter) (clientSession *services.ClientSession, err error) {
	clientSession, err = getClientSession(request)
	if err != nil {
		isNewClient := (err == http.ErrNoCookie)
		if isNewClient {
			fmt.Printf("Creating new client for %s\n", request.RemoteAddr)
			mutex.Lock()
			clientSequenceId++
			nextId := clientSequenceId
			mutex.Unlock()
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
