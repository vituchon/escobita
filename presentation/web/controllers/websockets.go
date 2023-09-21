package controllers

import (
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

func (h *WebSocketsHandler) Release(r *http.Request) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	clientId := h.clientIdResolverFunc(r)
	conn, exists := h.connByClientId[clientId]
	if !exists {
		return ConnectionDoesntExistErr
	}
	delete(h.connByClientId, clientId)
	_1000 := []byte{3, 232} // 1000, honouring https://tools.ietf.org/html/rfc6455#page-36
	conn.WriteMessage(websocket.CloseMessage, _1000)
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

func ReleaseWebSocket(w http.ResponseWriter, r *http.Request) {
	conn := webSocketsHandler.Retrieve(r)
	if conn != nil {
		gameWebSockets.UnbindClientWebSocketInGame(conn, r) // just in case the ws is associated with a game, then delete the association
		err := webSocketsHandler.Release(r)
		if err != nil {
			log.Printf("Error releasingWeb socket(RemoteAddr='%s') for client(id='%d')\n: %v\n", conn.RemoteAddr().String(), getWebPlayerId(r), err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			log.Printf("Web socket(RemoteAddr='%s') released OK for client(id='%d')\n", conn.RemoteAddr().String(), getWebPlayerId(r))
		}
	} else {
		log.Printf("No need to release web socket as it was not adquired (or already released) for  client(id='%d')\n", getWebPlayerId(r))
		w.WriteHeader(http.StatusBadRequest)
	}
}

func DebugWebSockets(w http.ResponseWriter, r *http.Request) {
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

	type response struct {
		Websockets      []websocket      `json:"websockets"`
		GamesWebsockets []gameWebsockets `json:"gamesWebsockets"`
	}
	WriteJsonResponse(w, http.StatusOK, response{Websockets: websockets, GamesWebsockets: gamesWebsockets})
}
