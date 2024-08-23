package services

import (
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ClientIdResolverFunc func(r *http.Request) int

type webSocketsHandler struct {
	upgrader             websocket.Upgrader
	ConnByClientId       map[int]*websocket.Conn
	mutex                sync.Mutex
	clientIdResolverFunc ClientIdResolverFunc
}

func NewWebSocketsHandler(clientIdResolverFunc ClientIdResolverFunc) webSocketsHandler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	connsByClientId := make(map[int]*websocket.Conn)
	return webSocketsHandler{
		upgrader:             upgrader,
		ConnByClientId:       connsByClientId,
		mutex:                sync.Mutex{},
		clientIdResolverFunc: clientIdResolverFunc,
	}
}

func (h *webSocketsHandler) AdquireOrRetrieve(w http.ResponseWriter, r *http.Request) (*websocket.Conn, bool, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	clientId := h.clientIdResolverFunc(r)
	conn, exists := h.ConnByClientId[clientId]
	if !exists { // server rules: at most one connection per http client (thus it is no multi tab compliant!)
		conn, err := h.upgrader.Upgrade(w, r, http.Header(map[string][]string{
			"created": []string{strconv.Itoa(int(time.Now().Unix()))},
		}))
		if err != nil {
			return nil, false, err
		}
		h.ConnByClientId[clientId] = conn
		return conn, true, nil // adquired (is new)
	}
	return conn, false, nil // retrieved (is not new)
}

func (h *webSocketsHandler) Retrieve(r *http.Request) *websocket.Conn {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	clientId := h.clientIdResolverFunc(r)
	return h.ConnByClientId[clientId]
}

func (h *webSocketsHandler) Release(r *http.Request, reason string) error {
	clientId := h.clientIdResolverFunc(r)
	return h.DoRelease(clientId, reason)
}

func (h *webSocketsHandler) DoRelease(clientId int, reason string) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	conn, exists := h.ConnByClientId[clientId]
	if !exists {
		return ConnectionDoesntExistErr
	}
	delete(h.ConnByClientId, clientId)
	//_1000 := []byte{3, 232}
	closeMessage := websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason) // honouring https://tools.ietf.org/html/rfc6455#page-36
	conn.WriteMessage(websocket.CloseMessage, closeMessage)
	return conn.Close()
}

var (
	ConnectionDoesntExistErr = errors.New("Connection doesn't exists")
	WebSocketsHandler        = NewWebSocketsHandler(GetClientPlayerId)
)
