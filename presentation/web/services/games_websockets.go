package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"sync"

	"github.com/gorilla/websocket"
	"github.com/vituchon/escobita/model"
	"github.com/vituchon/escobita/repositories"
)

type gameWebSockets struct {
	ConnsByGameId map[int][]*websocket.Conn
	mutex         sync.Mutex
}

var GameWebSockets gameWebSockets = gameWebSockets{ConnsByGameId: make(map[int][]*websocket.Conn)}

func (gws *gameWebSockets) NotifyGameConns(gameId int, kind string, data interface{}) {
	type Notification struct {
		Kind      string      `json:"kind"`
		BagOfCats interface{} `json:"data"`
	}

	gws.mutex.Lock()
	defer gws.mutex.Unlock()
	conns := gws.ConnsByGameId[gameId]
	for _, conn := range conns {
		notification := Notification{Kind: kind, BagOfCats: data}
		notificationAsJson, err := json.Marshal(notification)
		if err != nil {
			log.Printf("Error on marshalling notification, skip send. Error was: '%v'\n", err)
			continue
		}
		err = conn.WriteMessage(websocket.TextMessage, notificationAsJson)
		if err != nil {
			log.Println(err)
		}
	}
}

func (gws *gameWebSockets) BindClientWebSocketToGame(response http.ResponseWriter, request *http.Request, gameId int) {
	log.Printf("Binding web socket from client(id=%d) in game(id=%d)...", GetWebPlayerId(request), gameId)
	conn, _, err := WebSocketsHandler.AdquireOrRetrieve(response, request)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	gws.mutex.Lock()
	defer gws.mutex.Unlock()

	for _, exstingConn := range gws.ConnsByGameId[gameId] {
		if exstingConn == conn {
			msg := fmt.Sprintf("Web socket(remoteAddr='%s') from client(id=%d) already binded in game(id=%d)", conn.RemoteAddr().String(), GetWebPlayerId(request), gameId)
			log.Println(msg)
			http.Error(response, msg, http.StatusBadRequest)
			return
		}
	}

	gws.ConnsByGameId[gameId] = append(gws.ConnsByGameId[gameId], conn)
	log.Printf("Binded web socket(remoteAddr='%s') from client(id=%d) in game(id=%d)", conn.RemoteAddr().String(), GetWebPlayerId(request), gameId)
}

func (gws *gameWebSockets) UnbindAllWebSocketsInGame(gameId int, request *http.Request) {
	gws.mutex.Lock()
	defer gws.mutex.Unlock()
	log.Printf("Unbinding all web sockets from possible joined game id='%d'...\n", gameId)

	for _, conn := range gws.ConnsByGameId[gameId] {
		gws.doUnbindClientWebSocketInGame(conn, gameId, request)
	}
	delete(gws.ConnsByGameId, gameId)

	log.Printf("Unbinded all web sockets from possible joined game id='%d'\n", gameId)
}

func (gws *gameWebSockets) UnbindClientWebSocketInGame(conn *websocket.Conn, request *http.Request) {
	gws.mutex.Lock()
	defer gws.mutex.Unlock()
	log.Printf("Unbinding web socket(remoteAddr='%s') from a possible joined game...\n", conn.RemoteAddr().String())

	for gameId, conns := range gws.ConnsByGameId {
		for _, _conn := range conns {
			if _conn == conn {
				gws.doUnbindClientWebSocketInGame(conn, gameId, request)
				return
			}
		}
	}
	log.Printf("Web socket(remoteAddr='%s') was NOT binded to a game\n", conn.RemoteAddr().String())
}

// helper function, internal usage, do note that synchronization must be provided by in the client code... for now the only client is UnbindClientWebSocketInGame
func (gws *gameWebSockets) doUnbindClientWebSocketInGame(givenConn *websocket.Conn, gameId int, request *http.Request) {
	log.Printf("Unbinding Web socket(remoteAddr='%s') in game id='%d'...\n", givenConn.RemoteAddr().String(), gameId)
	conns := gws.ConnsByGameId[gameId]
	connsPtr := &conns
	chopped := (*connsPtr)[:0]
	for _, conn := range conns {
		if givenConn != conn {
			chopped = append(chopped, conn)
		}
	}
	*connsPtr = chopped
	gws.ConnsByGameId[gameId] = *connsPtr
	log.Printf("Unbinded Web socket(remoteAddr='%s') in game id='%d'\n", givenConn.RemoteAddr().String(), gameId)
}

type WebSockectOutgoingActionMsgPayload struct {
	Game   *repositories.PersistentGame `json:"game"`
	Action *model.PlayerAction          `json:"action,omitempty"`
}
