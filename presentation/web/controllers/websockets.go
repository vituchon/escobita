package controllers

import (
	"errors"
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

func (h *WebSocketsHandler) Adquire(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	clientId := h.clientIdResolverFunc(r)
	h.mutex.Lock()
	defer h.mutex.Unlock()
	conn, exists := h.connsByClientId[clientId]
	if !exists {
		conn, err := h.upgrader.Upgrade(w, r, http.Header(map[string][]string{
			"created": []string{strconv.Itoa(int(time.Now().Unix()))},
		}))
		if err != nil {
			return nil, err
		}
		h.connsByClientId[clientId] = conn
	}
	return conn, nil
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
	return conn.Close()
}

var (
	ConnectionDoesntExistErr = errors.New("Connection doesn't exists")
	webSocketsHandler        = NewWebSocketsHandler(getWebPlayerId)
)

func AdquireWebSocket(w http.ResponseWriter, r *http.Request) {
	_, err := webSocketsHandler.Adquire(w, r)
	if err != nil {
		log.Printf("Error getting web socket: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ReleaseWebSocket(w http.ResponseWriter, r *http.Request) {
	err := webSocketsHandler.Release(w, r)
	if err != nil {
		log.Printf("Error releasing web socket: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
