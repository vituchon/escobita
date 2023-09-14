package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/vituchon/escobita/model"
	"github.com/vituchon/escobita/presentation/web/services"
	"github.com/vituchon/escobita/repositories"

	"github.com/gorilla/websocket"
)

var gamesRepository repositories.Games = repositories.NewGamesMemoryRepository()

func GetGames(response http.ResponseWriter, request *http.Request) {
	games, err := gamesRepository.GetGames()
	if err != nil {
		log.Printf("error while retrieving games : '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, games)
}

func GetGameById(response http.ResponseWriter, request *http.Request) {
	paramId := RouteParam(request, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		log.Printf("Can not parse id from '%s'", paramId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	game, err := gamesRepository.GetGameById(id)
	if err != nil {
		log.Printf("error while retrieving game(id=%d): '%v'", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, game)
}

const MAX_GAMES_PER_PLAYER = 1

func CreateGame(response http.ResponseWriter, request *http.Request) {
	playerId := getWebPlayerId(request) // will be the game's owner
	if gamesRepository.GetGamesCreatedCount(playerId) == MAX_GAMES_PER_PLAYER {
		msg := fmt.Sprintf("Player(id='%d') has reached the maximum game creation limit: '%v'", playerId, MAX_GAMES_PER_PLAYER)
		response.WriteHeader(http.StatusBadRequest)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}

	var game repositories.PersistentGame
	err := parseJsonFromReader(request.Body, &game)
	if err != nil {
		log.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	game.PlayerId = playerId
	created, err := gamesRepository.CreateGame(game)
	if err != nil {
		log.Printf("error while creating Game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, created)
}

func UpdateGame(response http.ResponseWriter, request *http.Request) {
	var game repositories.PersistentGame
	err := parseJsonFromReader(request.Body, &game)
	if err != nil {
		log.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	updated, err := gamesRepository.UpdateGame(game)
	if err != nil {
		log.Printf("error while updating Game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	msgPayload := WebSockectOutgoingActionMsgPayload{updated, nil}
	gameWebSockets.NotifyGameConns(*game.Id, "updated", msgPayload)
	WriteJsonResponse(response, http.StatusOK, updated)
}

func DeleteGame(response http.ResponseWriter, request *http.Request) {
	paramId := RouteParam(request, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		log.Printf("Can not parse id from '%s'", paramId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	var player repositories.PersistentPlayer
	err = parseJsonFromReader(request.Body, &player)
	if err != nil {
		log.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	game, err := gamesRepository.GetGameById(id)
	if err != nil {
		log.Printf("error while retrieving game(id=%d): '%v'", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !services.CanPlayerDeleteGame(game, player) {
		log.Printf("Only game's owner(id=%d) is allowed to delete it. Requesting player(id='%d') is not the owner.", game.PlayerId, *player.Id)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	err = gamesRepository.DeleteGame(id)
	if err != nil {
		log.Printf("error while deleting game(id=%d): '%v'", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	gameWebSockets.UnbindAllWebSocketsInGame(id, request)
	response.WriteHeader(http.StatusOK)
}

// Escobita oriented events

func ResumeGame(response http.ResponseWriter, request *http.Request) {
	var game repositories.PersistentGame

	/*bufferedReader := bufio.NewReader(request.Body)

	// Read all data into a single buffer
	buffer, err := bufferedReader.ReadBytes(0) // 0 means to read until the end
	if err != nil && err != io.EOF {
		log.Printf("Error reading from reader: %v\n", err)
		return
	}

	// Print the entire content
	fmt.Println("Game:", string(buffer))

	err = parseJsonFromReader(bytes.NewReader(buffer), &game)*/
	err := parseJsonFromReader(request.Body, &game)
	if err != nil {
		log.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	playerId := getWebPlayerId(request)
	if game.PlayerId != playerId {
		log.Printf("error while starting Game: request doesn't cames from the owner, in cames from %d\n", playerId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	updated, err := services.ResumeGame(game)
	updated, err = gamesRepository.UpdateGame(*updated)
	if err != nil {
		log.Printf("error while starting Game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	msgPayload := WebSockectOutgoingActionMsgPayload{updated, nil}
	gameWebSockets.NotifyGameConns(*game.Id, "resume", msgPayload)
	WriteJsonResponse(response, http.StatusOK, updated)
}

type WebSockectOutgoingActionMsgPayload struct {
	Game   *repositories.PersistentGame `json:"game"`
	Action *model.PlayerAction          `json:"action,omitempty"`
}

func PerformTakeAction(response http.ResponseWriter, request *http.Request) {
	paramId := RouteParam(request, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		log.Printf("Can not parse id from '%s'", paramId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	game, err := gamesRepository.GetGameById(id)
	if err != nil {
		log.Printf("error getting game by id: '%v'", err)
		response.WriteHeader(http.StatusBadRequest) // dev note: it may be 404 NotFound is the case the game with the given id doesn't exists
		return
	}
	var takeAction model.PlayerTakeAction
	err = parseJsonFromReader(request.Body, &takeAction)
	if err != nil {
		log.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	updated, action := services.PerformTakeAction(*game, takeAction)
	updated, err = gamesRepository.UpdateGame(*updated)
	if err != nil {
		log.Printf("error while doing take action: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	msgPayload := WebSockectOutgoingActionMsgPayload{game, &action}
	gameWebSockets.NotifyGameConns(*game.Id, "take", msgPayload)
	WriteJsonResponse(response, http.StatusOK, msgPayload)
}

func PerformDropAction(response http.ResponseWriter, request *http.Request) {
	paramId := RouteParam(request, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		log.Printf("Can not parse id from '%s'", paramId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	game, err := gamesRepository.GetGameById(id)
	if err != nil {
		log.Printf("error getting game by id: '%v'", err)
		response.WriteHeader(http.StatusBadRequest) // dev note: it may be 404 NotFound is the case the game with the given id doesn't exists
		return
	}
	var dropAction model.PlayerDropAction
	err = parseJsonFromReader(request.Body, &dropAction)
	if err != nil {
		log.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	game, action := services.PerformDropAction(*game, dropAction)
	game, err = gamesRepository.UpdateGame(*game)
	if err != nil {
		log.Printf("error while doing drop action: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	msgPayload := WebSockectOutgoingActionMsgPayload{game, &action}
	gameWebSockets.NotifyGameConns(*game.Id, "drop", msgPayload)
	WriteJsonResponse(response, http.StatusOK, msgPayload)
}

func CalculateGameStats(response http.ResponseWriter, request *http.Request) {
	paramId := RouteParam(request, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		log.Printf("Can not parse id from '%s'", paramId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	paramMatchIndex, exists := request.URL.Query()["matchIndex"] // 0 is the first, 1 the second and so on...
	if !exists {
		log.Println("There was no matchIndex parameter in the query string")
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	matchIndex, err := strconv.Atoi(paramMatchIndex[0])
	if err != nil {
		log.Printf("Can not parse match index from '%s'", paramMatchIndex)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	game, err := gamesRepository.GetGameById(id)
	if err != nil {
		log.Printf("error getting game by id: '%v'", err)
		response.WriteHeader(http.StatusBadRequest) // dev note: it may be 404 NotFound is the case the game with the given id doesn't exists
		return
	}

	var stats model.ScoreSummaryByPlayer
	if matchIndex == len(game.Matchs) {
		stats = services.CalculateCurrentMatchStats(*game)
	} else {
		stats = services.CalculatePlayedMatchStats(*game, matchIndex)
	}

	if err != nil {
		log.Printf("error while calculating game stats action: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, stats)
}

type WebSockectOutgoingChatMsgPayload struct {
	Message services.VolatileWebMessage `json:"message"`
}

func SendMessage(response http.ResponseWriter, request *http.Request) {
	var message services.VolatileWebMessage
	err := parseJsonFromReader(request.Body, &message)
	if err != nil {
		log.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	paramId := RouteParam(request, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		log.Printf("Can not parse id from '%s'", paramId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	msgPayload := WebSockectOutgoingChatMsgPayload{message}
	gameWebSockets.NotifyGameConns(id, "game-chat", msgPayload)
	WriteJsonResponse(response, http.StatusOK, struct{}{})
}

type GameWebSockets struct {
	connsByGameId map[int][]*websocket.Conn
	mutex         sync.Mutex
}

var gameWebSockets GameWebSockets = GameWebSockets{connsByGameId: make(map[int][]*websocket.Conn)}

func (gws *GameWebSockets) NotifyGameConns(gameId int, kind string, data interface{}) {
	type Notification struct {
		Kind      string      `json:"kind"`
		BagOfCats interface{} `json:"data"`
	}

	gws.mutex.Lock()
	defer gws.mutex.Unlock()
	conns := gws.connsByGameId[gameId]
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

func BindClientWebSocketToGame(response http.ResponseWriter, request *http.Request) {
	gameId, err := ParseRouteParamAsInt(request, "id")
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	gameWebSockets.BindClientWebSocketToGame(response, request, gameId)
}

func (gws *GameWebSockets) BindClientWebSocketToGame(response http.ResponseWriter, request *http.Request, gameId int) {
	log.Printf("Binding web socket from client(id=%d) in game(id=%d)...", getWebPlayerId(request), gameId)
	conn, _, err := webSocketsHandler.AdquireOrRetrieve(response, request)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	gws.mutex.Lock()
	defer gws.mutex.Unlock()

	for _, exstingConn := range gws.connsByGameId[gameId] {
		if exstingConn == conn {
			log.Printf("Web socket(remoteAddr='%s') from client(id=%d) already binded in game(id=%d)", conn.RemoteAddr().String(), getWebPlayerId(request), gameId)
			return
		}
	}

	gws.connsByGameId[gameId] = append(gws.connsByGameId[gameId], conn)
	log.Printf("Binded web socket(remoteAddr='%s') from client(id=%d) in game(id=%d)", conn.RemoteAddr().String(), getWebPlayerId(request), gameId)
}

func UnbindClientWebSocketInGame(response http.ResponseWriter, request *http.Request) {
	conn := webSocketsHandler.Retrieve(request)
	if conn != nil {
		gameWebSockets.UnbindClientWebSocketInGame(conn, request)
		response.WriteHeader(http.StatusOK)
	} else {
		log.Printf("No need to release web socket as it was not adquired (or already released) for  client(id='%d')\n", getWebPlayerId(request))
		response.WriteHeader(http.StatusBadRequest)
	}
}

func (gws *GameWebSockets) UnbindAllWebSocketsInGame(gameId int, request *http.Request) {
	gws.mutex.Lock()
	defer gws.mutex.Unlock()
	log.Printf("Unbinding all web sockets from possible joined game id='%d'...\n", gameId)

	for _, conn := range gws.connsByGameId[gameId] {
		gws.doUnbindClientWebSocketInGame(conn, gameId, request)
	}
	delete(gws.connsByGameId, gameId)

	log.Printf("Unbinded all web sockets from possible joined game id='%d'\n", gameId)
}

func (gws *GameWebSockets) UnbindClientWebSocketInGame(conn *websocket.Conn, request *http.Request) {
	gws.mutex.Lock()
	defer gws.mutex.Unlock()
	log.Printf("Unbinding web socket(remoteAddr='%s') from possible joined game...\n", conn.RemoteAddr().String())

	for gameId, conns := range gws.connsByGameId {
		for _, _conn := range conns {
			if _conn == conn {
				gws.doUnbindClientWebSocketInGame(conn, gameId, request)
				return
			}
		}
	}
	log.Printf("Web socket(remoteAddr='%s') is NOT binded to a game\n", conn.RemoteAddr().String())
}

// helper function, internal usage, do note that synchronization must be provided by in the client code... for now the only client is UnbindClientWebSocketInGame
func (gws *GameWebSockets) doUnbindClientWebSocketInGame(givenConn *websocket.Conn, gameId int, request *http.Request) {
	log.Printf("Unbinding Web socket(remoteAddr='%s') in game id='%d'...\n", givenConn.RemoteAddr().String(), gameId)
	conns := gws.connsByGameId[gameId]
	connsPtr := &conns
	chopped := (*connsPtr)[:0]
	for _, conn := range conns {
		if givenConn != conn {
			chopped = append(chopped, conn)
		}
	}
	*connsPtr = chopped
	gws.connsByGameId[gameId] = *connsPtr
	log.Printf("Unbinded Web socket(remoteAddr='%s') in game id='%d'\n", givenConn.RemoteAddr().String(), gameId)
}
