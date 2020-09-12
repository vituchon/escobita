package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"local/escobita/model"
	"local/escobita/presentation/web/services"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"

	"github.com/gorilla/mux"
)

func Healthcheck(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
}

const ServerVersion = "0.0.1"

func Version(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte(ServerVersion))
	response.WriteHeader(http.StatusOK)
}

func WriteJsonResponse(response http.ResponseWriter, status int, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("error while mashalling object %v, trace: %+v", data, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_, err = response.Write(bytes)
	if err != nil {
		fmt.Printf("error while writting bytes to response writer: %+v", err)
	}
}

// Parses a json that cames from a reader an place in into the variable passed as argument
func ParseJsonFromReader(reader io.Reader, val interface{}) error {
	err := json.NewDecoder(reader).Decode(val)
	if err != nil {
		fmt.Printf("error decoding %T, error: %s", val, err.Error())
		return err
	}
	return nil
}

// Gets a Route parameter, that is a value within the url's PATH, not in the url's QUERY STRING.
func RouteParam(request *http.Request, name string) string {
	return mux.Vars(request)[name]
}

// WEB CLIENT SESSIONS

var clientSessions *sessions.CookieStore

func NewSessionStore(key []byte) {
	clientSessions = sessions.NewCookieStore(key)
}

var clientSequenceId int = 0

func GetOrCreateClientSession(request *http.Request) *sessions.Session {
	clientSession, err := clientSessions.Get(request, "client_session")
	if err != nil {
		fmt.Printf("error while retrieving 'client_session' from session store: %+v\n", err)
	}
	if clientSession.IsNew {
		fmt.Printf("creating new session\n")
		// begin not thread safe
		nextId := clientSequenceId + 1
		clientSession.Values["clientId"] = nextId
		clientSequenceId++
		// end not thread safe
	} else {
		fmt.Printf("using existing session, clientId = %v\n", clientSession.Values["clientId"])
	}
	return clientSession
}

func SaveClientSession(request *http.Request, response http.ResponseWriter, clientSession *sessions.Session) error {
	return clientSessions.Save(request, response, clientSession)
}

// GAMES

func GetGames(response http.ResponseWriter, request *http.Request) {
	games, err := services.GetGames()
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
	game, err := services.GetGameById(id)
	if err != nil {
		fmt.Printf("error while retrieving game(id=%d): '%v'", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, game)
}

func CreateGame(response http.ResponseWriter, request *http.Request) {
	var game services.WebGame
	err := ParseJsonFromReader(request.Body, &game)
	if err != nil {
		fmt.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	created, err := services.CreateGame(game)
	if err != nil {
		fmt.Printf("error while creating Game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("Tenemos este game creado %+v", created)
	WriteJsonResponse(response, http.StatusOK, created)
}

func UpdateGame(response http.ResponseWriter, request *http.Request) {
	var game services.WebGame
	err := ParseJsonFromReader(request.Body, &game)
	if err != nil {
		fmt.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	updated, err := services.UpdateGame(game)
	if err != nil {
		fmt.Printf("error while updating Game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
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
	err = services.DeleteGame(id)
	if err != nil {
		fmt.Printf("error while deleting game(id=%d): '%v'", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
}

func StartGame(response http.ResponseWriter, request *http.Request) {
	var game services.WebGame
	err := ParseJsonFromReader(request.Body, &game)
	if err != nil {
		fmt.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	updated, err := services.StartGame(game)
	if err != nil {
		fmt.Printf("error while starting Game: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, updated)
}

type gameActionResponseData struct {
	Game   *services.WebGame  `json:"game"`
	Action model.PlayerAction `json:"action"`
}

func PerformTakeAction(response http.ResponseWriter, request *http.Request) {
	paramId := RouteParam(request, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		fmt.Printf("Can not parse id from '%s'", paramId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	game, err := services.GetGameById(id)
	if err != nil {
		fmt.Printf("error getting game by id: '%v'", err)
		response.WriteHeader(http.StatusBadRequest) // dev note: it may be 404 NotFound is the case the game with the given id doesn't exists
		return
	}
	fmt.Printf("==========\ngame: %+v,\n============\n", *game)
	var takeAction model.PlayerTakeAction
	err = ParseJsonFromReader(request.Body, &takeAction)
	if err != nil {
		fmt.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	game, action, err := services.PerformTakeAction(*game, takeAction)
	if err != nil {
		fmt.Printf("error while doing take action: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, gameActionResponseData{game, action})
}

func PerformDropAction(response http.ResponseWriter, request *http.Request) {
	paramId := RouteParam(request, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		fmt.Printf("Can not parse id from '%s'", paramId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	game, err := services.GetGameById(id)
	if err != nil {
		fmt.Printf("error getting game by id: '%v'", err)
		response.WriteHeader(http.StatusBadRequest) // dev note: it may be 404 NotFound is the case the game with the given id doesn't exists
		return
	}
	var dropAction model.PlayerDropAction
	err = ParseJsonFromReader(request.Body, &dropAction)
	if err != nil {
		fmt.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	game, action, err := services.PerformDropAction(*game, dropAction)
	if err != nil {
		fmt.Printf("error while doing drop action: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, gameActionResponseData{game, action})
}

func CalculateGameStats(response http.ResponseWriter, request *http.Request) {
	paramId := RouteParam(request, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		fmt.Printf("Can not parse id from '%s'", paramId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	game, err := services.GetGameById(id)
	if err != nil {
		fmt.Printf("error getting game by id: '%v'", err)
		response.WriteHeader(http.StatusBadRequest) // dev note: it may be 404 NotFound is the case the game with the given id doesn't exists
		return
	}

	stats := services.CalculateGameStats(*game)
	if err != nil {
		fmt.Printf("error while calculating game stats action: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, stats)
}

// PLAYERS

func getWebPlayerId(request *http.Request) int {
	clientSession := request.Context().Value("clientSession").(*sessions.Session)
	wrappedInt, _ := clientSession.Values["clientId"]
	return wrappedInt.(int)
}

func ensurePlayerHasId(request *http.Request, player *services.WebPlayer) {
	if player.Id == nil {
		id := getWebPlayerId(request)
		player.Id = &id
	}
}

func GetPlayers(response http.ResponseWriter, request *http.Request) {
	players, err := services.GetPlayers()
	if err != nil {
		fmt.Printf("error while retrieving players : '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, players)
}

// Gets the web client's correspondant player
func GetClientPlayer(response http.ResponseWriter, request *http.Request) {
	id := getWebPlayerId(request)
	player, err := services.GetPlayerById(id)
	if err != nil {
		if err == services.EntityNotExistsErr {
			player = &services.WebPlayer{
				Player: model.Player{
					Name: "",
				},
				Id: &id,
			}
			player, err = services.CreatePlayer(*player)
			fmt.Printf("Creating new player %+v \n", player)
		}
		if err != nil {
			fmt.Printf("error while getting client player : '%v'", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		fmt.Printf("Using existing player %+v \n", player)
	}
	WriteJsonResponse(response, http.StatusOK, player)
}

func GetPlayerById(response http.ResponseWriter, request *http.Request) {
	paramId := RouteParam(request, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		fmt.Printf("Can not parse id from '%s'", paramId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	player, err := services.GetPlayerById(id)
	if err != nil {
		fmt.Printf("error while retrieving player(id=%d): '%v'\n", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, player)
}

func CreatePlayer(response http.ResponseWriter, request *http.Request) {
	var player services.WebPlayer
	err := ParseJsonFromReader(request.Body, &player)
	if err != nil {
		fmt.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("ParseJsonFromReader(request.Body, &player) = %v %v\n", player, err)
	ensurePlayerHasId(request, &player)
	fmt.Printf("ensurePlayerHasId(request, &player) => %v\n", player)

	created, err := services.CreatePlayer(player)
	if err != nil {
		fmt.Printf("error while creating Player: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, created)
}

func UpdatePlayer(response http.ResponseWriter, request *http.Request) {
	var player services.WebPlayer
	err := ParseJsonFromReader(request.Body, &player)
	if err != nil {
		fmt.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("ParseJsonFromReader(request.Body, &player) = %v %v\n", player, err)
	ensurePlayerHasId(request, &player)
	fmt.Printf("ensurePlayerHasId(request, &player) => %v\n", player)
	updated, err := services.UpdatePlayer(player)
	if err != nil {
		fmt.Printf("error while updating Player: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, updated)
}

func DeletePlayer(response http.ResponseWriter, request *http.Request) {
	paramId := RouteParam(request, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		fmt.Printf("Can not parse id from '%s'", paramId)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	err = services.DeletePlayer(id)
	if err != nil {
		fmt.Printf("error while deleting player(id=%d): '%v'", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
}
