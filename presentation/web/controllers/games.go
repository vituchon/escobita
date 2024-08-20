package controllers

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/vituchon/escobita/model"
	"github.com/vituchon/escobita/presentation/web/services"
	"github.com/vituchon/escobita/repositories"
)

// TODO: refact: promote usage of gameId in the following endpoints
/*	apiPost("/games/{id:[0-9]+}/message", controllers.SendMessage)
	//apiPut("/games/{id:[0-9]+}", controllers.UpdateGame)
	apiDelete("/games/{id:[0-9]+}", controllers.DeleteGame)
	apiPost("/games/{id:[0-9]+}/start", controllers.StartGame)
	apiPost("/games/{id:[0-9]+}/join", controllers.JoinGame)
	apiPost("/games/{id:[0-9]+}/quit", controllers.QuitGame)
	apiPost("/games/{id:[0-9]+}/perform-take-action", controllers.PerformTakeAction)
	apiPost("/games/{id:[0-9]+}/perform-drop-action", controllers.PerformDropAction)
	apiGet("/games/{id:[0-9]+}/calculate-stats", controllers.CalculateGameStats)*/
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
	id, err := ParseRouteParamAsInt(request, "id")
	if err != nil {
		log.Println(err)
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
	playerId := services.GetWebPlayerId(request) // will be the game's owner
	if gamesRepository.GetGamesCreatedCount(playerId) == MAX_GAMES_PER_PLAYER {
		msg := fmt.Sprintf("Player(id='%d') has reached the maximum game creation limit: '%v'", playerId, MAX_GAMES_PER_PLAYER)
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}

	game, err := retrieveGameByValue(request)
	if err != nil {
		log.Printf("error while retrieving game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	player, err := playersRepository.GetPlayerById(playerId)
	if err != nil {
		log.Printf("error getting player by id='%d': '%v'", playerId, err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	game.Owner = *player
	created, err := gamesRepository.CreateGame(*game)
	if err != nil {
		log.Printf("error while creating game: '%v'", err)
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
		log.Printf("error while updating game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	msgPayload := services.WebSockectOutgoingActionMsgPayload{updated, nil}
	services.GameWebSockets.NotifyGameConns(*game.Id, "updated", msgPayload)
	WriteJsonResponse(response, http.StatusOK, updated)
}

func DeleteGame(response http.ResponseWriter, request *http.Request) {
	var player repositories.PersistentPlayer
	err := parseJsonFromReader(request.Body, &player)
	if err != nil {
		log.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	game, err := retrieveGameByReference(request)
	if err != nil {
		log.Printf("error while retrieving game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !services.CanPlayerDeleteGame(game, player) {
		log.Printf("Only game's owner(id=%d) is allowed to delete it. Requesting player(id='%d') is not the owner.", game.Owner.Id, player.Id)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	id := *game.Id
	err = gamesRepository.DeleteGame(id)
	if err != nil {
		log.Printf("error while deleting game(id=%d): '%v'", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	services.GameWebSockets.UnbindAllWebSocketsInGame(id, request)
	response.WriteHeader(http.StatusOK)
}

// Escobita oriented events

func StartGame(response http.ResponseWriter, request *http.Request) {
	game, err := retrieveGameByReference(request)
	if err != nil {
		log.Printf("error while retrieving game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	playerId := services.GetWebPlayerId(request)
	if game.Owner.Id != playerId {
		log.Printf("error while starting game: request doesn't cames from the owner, in cames from %d\n", playerId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	updated, err := services.StartGame(*game)
	updated, err = gamesRepository.UpdateGame(*updated)
	if err != nil {
		log.Printf("error while starting game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	msgPayload := services.WebSockectOutgoingActionMsgPayload{updated, nil}
	services.GameWebSockets.NotifyGameConns(*game.Id, "start", msgPayload)
	WriteJsonResponse(response, http.StatusOK, updated)
}

func JoinGame(response http.ResponseWriter, request *http.Request) {
	game, err := retrieveGameByReference(request)
	if err != nil {
		log.Printf("error while retrieving game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	playerId := services.GetWebPlayerId(request)
	player, err := playersRepository.GetPlayerById(playerId)
	if err != nil {
		msg := fmt.Sprintf("error while getting player by id, error was: '%v'\n", player)
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}

	err = game.Join(*player)
	if err != nil {
		msg := fmt.Sprintf("error while joining game, error was: '%v'\n", err)
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}
	updated, err := gamesRepository.UpdateGame(*game)
	if err != nil {
		log.Printf("error while updating game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	msgPayload := services.WebSockectOutgoingJoinMsgPayload{updated, player}
	services.GameWebSockets.NotifyGameConns(*game.Id, "join", msgPayload)
	WriteJsonResponse(response, http.StatusOK, game)
}

func QuitGame(response http.ResponseWriter, request *http.Request) {
	game, err := retrieveGameByReference(request)
	if err != nil {
		log.Printf("error while retrieving game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	playerId := services.GetWebPlayerId(request)
	player, err := playersRepository.GetPlayerById(playerId)
	if err != nil {
		msg := fmt.Sprintf("error while getting player by id, error was: '%v'\n", player)
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}

	err = game.Quit(*player)
	if err != nil {
		msg := fmt.Sprintf("error while quiting game, error was: '%v'\n", player)
		log.Println(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}
	updated, err := gamesRepository.UpdateGame(*game)
	if err != nil {
		log.Printf("error while updating game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	msgPayload := services.WebSockectOutgoingJoinMsgPayload{updated, player}
	services.GameWebSockets.NotifyGameConns(*game.Id, "quit", msgPayload)
	WriteJsonResponse(response, http.StatusOK, game)
}

func AddComputer(response http.ResponseWriter, request *http.Request) {
	game, err := retrieveGameByReference(request)
	if err != nil {
		log.Printf("error while retrieving game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	game.Join(model.ComputerPlayer)
	updated, err := gamesRepository.UpdateGame(*game)
	if err != nil {
		log.Printf("error while updating game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	message := services.VolatileWebMessage{Player: model.ComputerPlayer, Text: "Buenas y santas! 🧉😸"}
	chatMsgPayload := WebSockectOutgoingChatMsgPayload{message}
	services.GameWebSockets.NotifyGameConns(*game.Id, "game-chat", chatMsgPayload)

	joinMsgPayload := services.WebSockectOutgoingJoinMsgPayload{updated, &model.ComputerPlayer}
	services.GameWebSockets.NotifyGameConns(*game.Id, "join", joinMsgPayload)
	WriteJsonResponse(response, http.StatusOK, updated)
}

func PerformTakeAction(response http.ResponseWriter, request *http.Request) {
	game, err := retrieveGameByReference(request)
	if err != nil {
		log.Printf("error while retrieving game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return

	}
	var takeAction model.PlayerTakeAction
	err = parseJsonFromReader(request.Body, &takeAction)
	if err != nil {
		log.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	updated, action, err := services.PerformTakeAction(*game, takeAction)
	if err != nil {
		log.Printf("error while performing take action: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	updated, err = gamesRepository.UpdateGame(*updated)
	if err != nil {
		log.Printf("error while updating game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	msgPayload := services.WebSockectOutgoingActionMsgPayload{game, action}
	services.GameWebSockets.NotifyGameConns(*game.Id, "take", msgPayload)
	WriteJsonResponse(response, http.StatusOK, msgPayload)
}

func PerformDropAction(response http.ResponseWriter, request *http.Request) {
	game, err := retrieveGameByReference(request)
	if err != nil {
		log.Printf("error while retrieving game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return

	}
	var dropAction model.PlayerDropAction
	err = parseJsonFromReader(request.Body, &dropAction)
	if err != nil {
		log.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	game, action, err := services.PerformDropAction(*game, dropAction)
	if err != nil {
		log.Printf("error while performing drop action: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	game, err = gamesRepository.UpdateGame(*game)
	if err != nil {
		log.Printf("error while updating game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	msgPayload := services.WebSockectOutgoingActionMsgPayload{game, action}

	services.GameWebSockets.NotifyGameConns(*game.Id, "drop", msgPayload)
	WriteJsonResponse(response, http.StatusOK, msgPayload)
}

func CalculateGameStats(response http.ResponseWriter, request *http.Request) {

	matchIndex, err := ParseSingleIntegerUrlQueryParam(request, "matchIndex") // 0 is the first, 1 the second and so on...
	if err != nil {
		log.Printf("Can not parse match index: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	game, err := retrieveGameByReference(request)
	if err != nil {
		log.Printf("error while retrieving game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return

	}

	var stats model.ScoreSummaryByPlayer
	if *matchIndex == len(game.Matchs) {
		stats = services.CalculateCurrentMatchStats(*game)
	} else {
		stats = services.CalculatePlayedMatchStats(*game, *matchIndex)
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

	game, err := retrieveGameByReference(request)
	if err != nil {
		log.Printf("Error retrieving game for sending an obsequent accompaniment to vitu: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return

	}
	id := *game.Id
	msgPayload := WebSockectOutgoingChatMsgPayload{message}
	services.GameWebSockets.NotifyGameConns(id, "game-chat", msgPayload)

	if strings.Contains(strings.ToLower(message.Player.Name), "vitu") {
		if game.IsJoined(model.ComputerPlayer) {
			mustSendAbsequentAccompaniment := rand.Intn(2) == 0 // flip a coin
			if mustSendAbsequentAccompaniment {
				sendAbsequentAccompaniment := func() {
					names := []string{}
					for _, gamePlayer := range game.Players {
						if gamePlayer.Id != message.Player.Id && gamePlayer.Id != model.ComputerPlayer.Id {
							names = append(names, gamePlayer.Name)
						}
					}
					allNames := strings.Join(names, " ")
					var obsequentAccompaniments []string = []string{
						"Totalmente",
						"Kpo++ vituchon 👏🏿🤗",
						"Cuanta razón....",
						allNames + " calentito los panchos",
						allNames + " y todos los que lo leen son de la B",
						allNames + " silenció.... habló el maestro...",
						"Faa como la cantó...",
						allNames + " i-m-p-e-c-a-b-l-e lo que dijo el vitul 👏🏿",
						allNames + " shhh no saben nada...",
					}

					message.Player = model.ComputerPlayer
					message.Text = obsequentAccompaniments[rand.Intn(len(obsequentAccompaniments))]
					msgPayload := WebSockectOutgoingChatMsgPayload{message}
					services.GameWebSockets.NotifyGameConns(id, "game-chat", msgPayload)
				}
				time.AfterFunc(1*time.Second, sendAbsequentAccompaniment)
			}
		}

	}

	WriteJsonResponse(response, http.StatusOK, struct{}{})
}

func BindClientWebSocketToGame(response http.ResponseWriter, request *http.Request) {
	gameId, err := ParseRouteParamAsInt(request, "id")
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	services.GameWebSockets.BindClientWebSocketToGame(response, request, gameId)
}

func UnbindClientWebSocketInGame(response http.ResponseWriter, request *http.Request) {
	conn := services.WebSocketsHandler.Retrieve(request)
	if conn != nil {
		services.GameWebSockets.UnbindClientWebSocketInGame(conn, request)
		response.WriteHeader(http.StatusOK)
	} else {
		log.Printf("No need to release web socket as it was not adquired (or already released) for  client(id='%d')\n", services.GetWebPlayerId(request))
		response.WriteHeader(http.StatusBadRequest)
	}
}

// Retrieves the stored game in the underlying storage system using the id present in the URL (route param)
func retrieveGameByReference(request *http.Request) (*repositories.PersistentGame, error) {
	id, err := ParseRouteParamAsInt(request, "id")
	if err != nil {
		return nil, err
	}

	game, err := gamesRepository.GetGameById(id)
	if err != nil {
		errMsg := fmt.Sprintf("error while retrieving game(id=%d): '%v'", id, err)
		return nil, errors.New(errMsg)
	}
	return game, nil
}

// Retrieves the stored game in the request's payload
func retrieveGameByValue(request *http.Request) (*repositories.PersistentGame, error) {
	var game repositories.PersistentGame
	/*bufferedReader := bufio.NewReader(request.Body)

	// Read all data into a single buffer
	buffer, err := bufferedReader.ReadBytes(0) // 0 means to read until the end
	if err != nil && err != io.EOF {
		errMsg := fmt.Sprintf("Error reading from reader: %v\n", er)
		return nil, errors.New(errMsg)
	}

	// Print the entire content
	fmt.Println("Game:", string(buffer))

	err = parseJsonFromReader(bytes.NewReader(buffer), &game)*/

	err := parseJsonFromReader(request.Body, &game)
	if err != nil {
		errMsg := fmt.Sprintf("error reading request body: '%v'", err)
		return nil, errors.New(errMsg)
	}
	return &game, nil
}
