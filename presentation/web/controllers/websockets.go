package controllers

import (
	"errors"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ClientIdResolverFunc func(r *http.Request) int

type WebSocketsHandler struct {
	upgrader             websocket.Upgrader
	connsByClientId      map[int]*websocket.Conn
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
		connsByClientId:      connsByClientId,
		mutex:                sync.Mutex{},
		clientIdResolverFunc: clientIdResolverFunc,
	}
}

func (h *WebSocketsHandler) AdquireOrRetrieve(w http.ResponseWriter, r *http.Request) (*websocket.Conn, bool, error) {
	clientId := h.clientIdResolverFunc(r)
	h.mutex.Lock()
	defer h.mutex.Unlock()
	conn, exists := h.connsByClientId[clientId]
	if !exists { // server rules: at most one connection per http client (thus it is no multi tab compliant!)
		conn, err := h.upgrader.Upgrade(w, r, http.Header(map[string][]string{
			"created": []string{strconv.Itoa(int(time.Now().Unix()))},
		}))
		if err != nil {
			return nil, false, err
		}
		h.connsByClientId[clientId] = conn
		return conn, true, nil // adquired (is new)
	}
	return conn, false, nil // retrieved (is not new)
}

func (h *WebSocketsHandler) Release(w http.ResponseWriter, r *http.Request) error {
	clientId := h.clientIdResolverFunc(r)
	h.mutex.Lock()
	defer h.mutex.Unlock()
	conn, exists := h.connsByClientId[clientId]
	if !exists {
		return ConnectionDoesntExistErr
	}
	delete(h.connsByClientId, clientId)
	unbindSocketSocketFromJoinedGame(conn, r)
	_1000 := []byte{3, 232} // 1000, honouring https://tools.ietf.org/html/rfc6455#page-36
	conn.WriteMessage(websocket.CloseMessage, _1000)
	return conn.Close()
}

// for the sake of doing some sweep of possible binded websocket to a given game whose connection is about to be closed
func unbindSocketSocketFromJoinedGame(conn *websocket.Conn, request *http.Request) {
	for gameId, wss := range wsByGameId {
		for _, ws := range wss {
			if ws == conn {
				log.Println("ws is binded to a game, procedding to unbind...")
				UnbindConn(conn, gameId, request)
				return
			}
		}
	}
}

var (
	ConnectionDoesntExistErr = errors.New("Connection doesn't exists")
	webSocketsHandler        = NewWebSocketsHandler(getWebPlayerId)
)

func AdquireWebSocket(w http.ResponseWriter, r *http.Request) {
	_, isNew, err := webSocketsHandler.AdquireOrRetrieve(w, r)
	if err != nil {
		log.Printf("Error getting web socket: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	if !isNew {
		log.Printf("Web socket already adquired for client(id='%d')\n", getWebPlayerId(r))
		w.WriteHeader(http.StatusBadRequest)
	} else {
		log.Printf("Web socket adquired OK (connection \"upgraded\") for client(id='%d')\n", getWebPlayerId(r))
	}
}

func ReleaseWebSocket(w http.ResponseWriter, r *http.Request) {
	err := webSocketsHandler.Release(w, r)
	if err != nil {
		log.Printf("Error releasing web socket: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		log.Printf("Web socket released OK for client(id='%d')\n", getWebPlayerId(r))
	}
}

func DebugWebSockets(w http.ResponseWriter, r *http.Request) {
	conns := webSocketsHandler.connsByClientId
	type item struct {
		RemoteAddr net.Addr `json:"remoteAddr"`
	}
	items := []item{}
	for _, conn := range conns {
		items = append(items, item{RemoteAddr: conn.RemoteAddr()})
	}
	WriteJsonResponse(w, http.StatusOK, items)
}
