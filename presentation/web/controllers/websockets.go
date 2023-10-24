package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ClientIdResolverFunc func(r *http.Request) int

type WebSocketsHandler struct {
	upgrader             websocket.Upgrader
	connByClientId       map[int]*websocket.Conn
	mutex                sync.Mutex
	clientIdResolverFunc ClientIdResolverFunc
}

func NewWebSocketsHandler(clientIdResolverFunc ClientIdResolverFunc) WebSocketsHandler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	connsByClientId := make(map[int]*websocket.Conn)
	return WebSocketsHandler{
		upgrader:             upgrader,
		connByClientId:       connsByClientId,
		mutex:                sync.Mutex{},
		clientIdResolverFunc: clientIdResolverFunc,
	}
}

func (h *WebSocketsHandler) AdquireOrRetrieve(w http.ResponseWriter, r *http.Request) (*websocket.Conn, bool, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	clientId := h.clientIdResolverFunc(r)
	conn, exists := h.connByClientId[clientId]
	if !exists { // server rules: at most one connection per http client (thus it is no multi tab compliant!)
		conn, err := h.upgrader.Upgrade(w, r, http.Header(map[string][]string{
			"created": []string{strconv.Itoa(int(time.Now().Unix()))},
		}))
		if err != nil {
			return nil, false, err
		}
		h.connByClientId[clientId] = conn
		return conn, true, nil // adquired (is new)
	}
	return conn, false, nil // retrieved (is not new)
}

func (h *WebSocketsHandler) Retrieve(r *http.Request) *websocket.Conn {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	clientId := h.clientIdResolverFunc(r)
	return h.connByClientId[clientId]
}

func (h *WebSocketsHandler) Release(r *http.Request, reason string) error {
	clientId := h.clientIdResolverFunc(r)
	return h.doRelease(clientId, reason)
}

func (h *WebSocketsHandler) doRelease(clientId int, reason string) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	conn, exists := h.connByClientId[clientId]
	if !exists {
		return ConnectionDoesntExistErr
	}
	delete(h.connByClientId, clientId)
	//_1000 := []byte{3, 232}
	closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason) // honouring https://tools.ietf.org/html/rfc6455#page-36
	conn.WriteMessage(websocket.CloseMessage, closeMessage)
	return conn.Close()
}

var (
	ConnectionDoesntExistErr = errors.New("Connection doesn't exists")
	webSocketsHandler        = NewWebSocketsHandler(getWebPlayerId)
)

// WEB SOCKET DEDICATED END POINTS

func AdquireWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, isNew, err := webSocketsHandler.AdquireOrRetrieve(w, r)
	if err != nil {
		log.Printf("Error adquiring or retrieving web socket: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	if !isNew {
		msg := fmt.Sprintf("Web socket already adquired for client(id='%d')", getWebPlayerId(r))
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
	} else {
		log.Printf("Web socket(RemoteAddr='%s') adquired OK (connection \"upgraded\") for client(id='%d')\n", conn.RemoteAddr().String(), getWebPlayerId(r))
	}
}

func ReleaseWebSocket(response http.ResponseWriter, request *http.Request) {
	conn := webSocketsHandler.Retrieve(request)
	if conn != nil {
		gameWebSockets.UnbindClientWebSocketInGame(conn, request) // just in case the ws is associated with a game, then delete the association
		err := webSocketsHandler.Release(request, "Connection closed gracefully")
		if err != nil {
			log.Printf("Error releasingWeb socket(RemoteAddr='%s') for client(id='%d')\n: %v\n", conn.RemoteAddr().String(), getWebPlayerId(request), err)
			response.WriteHeader(http.StatusInternalServerError)
		} else {
			log.Printf("Web socket(RemoteAddr='%s') released OK for client(id='%d')\n", conn.RemoteAddr().String(), getWebPlayerId(request))
		}
	} else {
		log.Printf("No need to release web socket as it was not adquired (or already released) for  client(id='%d')\n", getWebPlayerId(request))
		response.WriteHeader(http.StatusBadRequest)
	}
}

type ServerMessage struct {
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

func SendMessageWebSocket(response http.ResponseWriter, request *http.Request) {
	playerId := getWebPlayerId(request)
	message, err := getPingMessageOrDefault(request)
	if err != nil {
		log.Printf("error getting ping message: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	conn := webSocketsHandler.connByClientId[playerId]
	if conn == nil {
		log.Printf("There is no web socket for client(id='%d')\n", playerId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	serverMsg := ServerMessage{Kind: "debug", Message: *message}
	serverMsgAsJson, err := json.Marshal(serverMsg)
	if err != nil {
		log.Printf("Error on marshalling server message, skip send. Error was: '%v'\n", err)
	}
	err = conn.WriteMessage(websocket.TextMessage, serverMsgAsJson)
	if err != nil {
		log.Printf("Web socket(RemoteAddr='%s') for client(id='%d') could not be used to send message, error was: '%v'\n", conn.RemoteAddr().String(), playerId, err)
	}
	response.WriteHeader(http.StatusOK)
}

func SendMessageAllWebSockets(response http.ResponseWriter, request *http.Request) {
	message, err := getPingMessageOrDefault(request)
	if err != nil {
		log.Printf("error getting ping message: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	for clientId, conn := range webSocketsHandler.connByClientId {
		serverMsg := ServerMessage{Kind: "debug", Message: *message}
		serverMsgAsJson, err := json.Marshal(serverMsg)
		if err != nil {
			log.Printf("Error on marshalling server message, skip send. Error was: '%v'\n", err)
		}
		err = conn.WriteMessage(websocket.TextMessage, serverMsgAsJson)
		if err != nil {
			log.Printf("Web socket(RemoteAddr='%s') for client(id='%d') could not be used to send message, error was: '%v'\n", conn.RemoteAddr().String(), clientId, err)
		}
	}
	response.WriteHeader(http.StatusOK)
}

func getPingMessageOrDefault(request *http.Request) (*string, error) {
	var defaultPingMessage string = "ping"

	message, err := ParseSingleStringUrlQueryParam(request, "message")
	if err != nil {
		if err == UrlQueryParamNotFoundErr {
			message = &defaultPingMessage
		} else {
			return nil, err
		}
	}
	return message, nil
}

func DebugWebSockets(response http.ResponseWriter, request *http.Request) {
	type websocket struct {
		ClientId   int    `json:"clientId"`
		RemoteAddr string `json:"remoteAddr"`
	}
	websockets := []websocket{}
	for clientId, conn := range webSocketsHandler.connByClientId {
		websockets = append(websockets, websocket{ClientId: clientId, RemoteAddr: conn.RemoteAddr().String()})
	}

	type gameWebsockets struct {
		GameId      int      `json:"gameId"`
		RemoteAddrs []string `json:"RemoteAddrs"`
	}
	gamesWebsockets := []gameWebsockets{}
	for gameId, conns := range gameWebSockets.connsByGameId {
		var remoteAddrs []string
		for _, conn := range conns {
			remoteAddrs = append(remoteAddrs, conn.RemoteAddr().String())
		}
		gamesWebsockets = append(gamesWebsockets, gameWebsockets{GameId: gameId, RemoteAddrs: remoteAddrs})
	}

	type data struct {
		Websockets      []websocket      `json:"websockets"`
		GamesWebsockets []gameWebsockets `json:"gamesWebsockets"`
	}
	WriteJsonResponse(response, http.StatusOK, data{Websockets: websockets, GamesWebsockets: gamesWebsockets})
}

func ReleaseBrokenWebSockets(response http.ResponseWriter, request *http.Request) {
	for clientId, conn := range webSocketsHandler.connByClientId {
		// TODO: research about it using “ping/pong frames” are used to check the connection, sent from the server, the browser responds to these automatically. See https://javascript.info/websocket#data-transfer
		// right now using a "home made recipe": trying twice: from https://gosamples.dev/broken-pipe/ it states that "The first write to the closed connection causes the peer to reply with an RST packet indicating that the connection should be terminated immediately. The second write to the socket that has already received the RST causes the broken pipe error"
		err1 := conn.WriteMessage(websocket.TextMessage, []byte(""))
		err2 := conn.WriteMessage(websocket.TextMessage, []byte(""))
		if err1 != nil || err2 != nil {
			log.Printf("Web socket(RemoteAddr='%s') for client(id='%d') appears to be broken (err1='%v',err2='%v'), releasing...!\n", conn.RemoteAddr().String(), clientId, err1, err2)
			releaseErr := webSocketsHandler.doRelease(clientId, "Connection appears to be broken")
			if releaseErr != nil {
				log.Printf("Error while releasing web socket(RemoteAddr='%s') for client(id='%d'): '%v'", conn.RemoteAddr().String(), clientId, releaseErr)
			} else {
				log.Printf("Released web socket(RemoteAddr='%s') for client(id='%d') gracefully", conn.RemoteAddr().String(), clientId)
			}
		}
	}
	response.WriteHeader(http.StatusOK)
}

func ReleaseAllWebSockets(response http.ResponseWriter, request *http.Request) {
	for clientId, conn := range webSocketsHandler.connByClientId {
		err := webSocketsHandler.doRelease(clientId, "Connection terminated by force")
		if err != nil {
			log.Printf("Error while releasing web socket(RemoteAddr='%s') for client(id='%d'): '%v'", conn.RemoteAddr().String(), clientId, err)
		} else {
			log.Printf("Released web socket(RemoteAddr='%s') for client(id='%d') gracefully", conn.RemoteAddr().String(), clientId)
		}
	}
	response.WriteHeader(http.StatusOK)
}
