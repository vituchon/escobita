package services

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func GetClientId(request *http.Request) int {
	clientSession := request.Context().Value("clientSession").(*sessions.Session)
	wrappedInt, _ := clientSession.Values["clientId"]
	return wrappedInt.(int)
}
