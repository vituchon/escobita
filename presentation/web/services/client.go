package services

import (
	//"log"
	"net/http"
)

type ClientSession struct {
	Id int
}

func GetClientPlayerId(request *http.Request) int {
	clientSession := request.Context().Value("clientSession").(*ClientSession)
	//log.Printf("For request ip %s got %+v client session", request.RemoteAddr, *clientSession)
	return clientSession.Id
}
