package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/vituchon/escobita/model"
	"github.com/vituchon/escobita/presentation/web/services"
	"github.com/vituchon/escobita/repositories"

	"github.com/gorilla/websocket"
)

var gamesRepository repositories.Games = repositories.NewGamesMemoryStorage()

func GetGames(response http.ResponseWriter, request *http.Request) {
	games, err := gamesRepository.GetGames()
	if err != nil {
		fmt.Printf("error while retrieving games : '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, games)
}

func GetGameById(response http.ResponseWriter, request *http.Request) {
	paramId := RouteParam(request, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		fmt.Printf("Can not parse id from '%s'", paramId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	game, err := gamesRepository.GetGameById(id)
	if err != nil {
		fmt.Printf("error while retrieving game(id=%d): '%v'", id, err)
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
		fmt.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	game.PlayerId = playerId
	created, err := gamesRepository.CreateGame(game)
	if err != nil {
		fmt.Printf("error while creating Game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, created)
}

func UpdateGame(response http.ResponseWriter, request *http.Request) {
	var game repositories.PersistentGame
	err := parseJsonFromReader(request.Body, &game)
	if err != nil {
		fmt.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	updated, err := gamesRepository.UpdateGame(game)
	if err != nil {
		fmt.Printf("error while updating Game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	msgPayload := WebSockectOutgoingMsgPayload{updated, nil}
	notifyBindedWebSockets(*game.Id, "updated", msgPayload)
	WriteJsonResponse(response, http.StatusOK, updated)
}

func DeleteGame(response http.ResponseWriter, request *http.Request) {
	paramId := RouteParam(request, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		fmt.Printf("Can not parse id from '%s'", paramId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	var player repositories.PersistentPlayer
	err = parseJsonFromReader(request.Body, &player)
	if err != nil {
		fmt.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	game, err := gamesRepository.GetGameById(id)
	if err != nil {
		fmt.Printf("error while retrieving game(id=%d): '%v'", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !services.CanPlayerDeleteGame(game, player) {
		fmt.Printf("Only game's owner(id=%d) is allowed to delete it. Requesting player(id='%d') is not the owner.", game.PlayerId, *player.Id)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	err = gamesRepository.DeleteGame(id)
	if err != nil {
		fmt.Printf("error while deleting game(id=%d): '%v'", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
}

// Escobita oriented events

func ResumeGame(response http.ResponseWriter, request *http.Request) {
	var game repositories.PersistentGame

	/*bufferedReader := bufio.NewReader(request.Body)

	// Read all data into a single buffer
	buffer, err := bufferedReader.ReadBytes(0) // 0 means to read until the end
	if err != nil && err != io.EOF {
		fmt.Printf("Error reading from reader: %v\n", err)
		return
	}

	// Print the entire content
	fmt.Println("Game:", string(buffer))

	err = parseJsonFromReader(bytes.NewReader(buffer), &game)*/
	err := parseJsonFromReader(request.Body, &game)
	if err != nil {
		fmt.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	playerId := getWebPlayerId(request)
	if game.PlayerId != playerId {
		fmt.Printf("error while starting Game: request doesn't cames from the owner, in cames from %d\n", playerId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	updated, err := services.ResumeGame(game)
	updated, err = gamesRepository.UpdateGame(*updated)
	if err != nil {
		fmt.Printf("error while starting Game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	msgPayload := WebSockectOutgoingMsgPayload{updated, nil}
	notifyBindedWebSockets(*game.Id, "resume", msgPayload)
	WriteJsonResponse(response, http.StatusOK, updated)
}

type WebSockectOutgoingMsgPayload struct {
	Game   *repositories.PersistentGame `json:"game"`
	Action *model.PlayerAction          `json:"action,omitempty"`
}

func PerformTakeAction(response http.ResponseWriter, request *http.Request) {
	paramId := RouteParam(request, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		fmt.Printf("Can not parse id from '%s'", paramId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	game, err := gamesRepository.GetGameById(id)
	if err != nil {
		fmt.Printf("error getting game by id: '%v'", err)
		response.WriteHeader(http.StatusBadRequest) // dev note: it may be 404 NotFound is the case the game with the given id doesn't exists
		return
	}
	fmt.Printf("==========\ngame: %+v,\n============\n", *game)
	var takeAction model.PlayerTakeAction
	err = parseJsonFromReader(request.Body, &takeAction)
	if err != nil {
		fmt.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	updated, action := services.PerformTakeAction(*game, takeAction)
	updated, err = gamesRepository.UpdateGame(*updated)
	if err != nil {
		fmt.Printf("error while doing take action: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	msgPayload := WebSockectOutgoingMsgPayload{game, &action}
	notifyBindedWebSockets(*game.Id, "take", msgPayload)
	WriteJsonResponse(response, http.StatusOK, msgPayload)
}

func PerformDropAction(response http.ResponseWriter, request *http.Request) {
	paramId := RouteParam(request, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		fmt.Printf("Can not parse id from '%s'", paramId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	game, err := gamesRepository.GetGameById(id)
	if err != nil {
		fmt.Printf("error getting game by id: '%v'", err)
		response.WriteHeader(http.StatusBadRequest) // dev note: it may be 404 NotFound is the case the game with the given id doesn't exists
		return
	}
	var dropAction model.PlayerDropAction
	err = parseJsonFromReader(request.Body, &dropAction)
	if err != nil {
		fmt.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	game, action := services.PerformDropAction(*game, dropAction)
	game, err = gamesRepository.UpdateGame(*game)
	if err != nil {
		fmt.Printf("error while doing drop action: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	msgPayload := WebSockectOutgoingMsgPayload{game, &action}
	notifyBindedWebSockets(*game.Id, "drop", msgPayload)
	WriteJsonResponse(response, http.StatusOK, msgPayload)
}

func CalculateGameStats(response http.ResponseWriter, request *http.Request) {
	paramId := RouteParam(request, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		fmt.Printf("Can not parse id from '%s'", paramId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	paramMatchIndex, exists := request.URL.Query()["matchIndex"] // 0 is the first, 1 the second and so on...
	if !exists {
		fmt.Println("There was no matchIndex parameter in the query string")
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	matchIndex, err := strconv.Atoi(paramMatchIndex[0])
	if err != nil {
		fmt.Printf("Can not parse match index from '%s'", paramMatchIndex)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	game, err := gamesRepository.GetGameById(id)
	if err != nil {
		fmt.Printf("error getting game by id: '%v'", err)
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
		fmt.Printf("error while calculating game stats action: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, stats)
}

var wsByGameId map[int][]*websocket.Conn = make(map[int][]*websocket.Conn) // TODO : Access may be has to be monitored by sync locks!

func notifyBindedWebSockets(gameId int, kind string, data interface{}) {
	type Notification struct {
		Kind      string      `json:"kind"`
		BagOfCats interface{} `json:"data"`
	}

	log.Printf("Notify clients about event(type=%s) ", kind)
	conns := wsByGameId[gameId]
	for _, conn := range conns {
		notification := Notification{Kind: kind, BagOfCats: data}
		notificationAsJson, err := json.Marshal(notification)
		if err != nil {
			log.Println(err, " NO SEND ")
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
	conn, _, err := webSocketsHandler.AdquireOrRetrieve(response, request)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	wsByGameId[gameId] = append(wsByGameId[gameId], conn)
	log.Printf("Bind ws from client(id=%d) into game(id=%d) using conn=%v", getWebPlayerId(request), gameId, conn.RemoteAddr())
}

func UnbindClientWebSocketToGame(response http.ResponseWriter, request *http.Request) {
	gameId, err := ParseRouteParamAsInt(request, "id")
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	clientConn, isNew, err := webSocketsHandler.AdquireOrRetrieve(response, request)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	if isNew {
		log.Println("Suspicious call to UnbindClientWebSocketToGame! web socket connection wasn't established before!")
	}
	UnbindConn(clientConn, gameId, request)
}

func UnbindConn(givenConn *websocket.Conn, gameId int, request *http.Request) {
	conns := wsByGameId[gameId]
	connsPtr := &conns
	chopped := (*connsPtr)[:0]
	for _, conn := range conns {
		if givenConn != conn {
			chopped = append(chopped, conn)
		}
	}
	*connsPtr = chopped
	wsByGameId[gameId] = *connsPtr
	log.Printf("Unbind ws from client(id=%d) for game(id=%d). conn was=%v", getWebPlayerId(request), gameId, givenConn.RemoteAddr())
}
