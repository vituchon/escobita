package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/vituchon/escobita/model"
	"github.com/vituchon/escobita/repositories"

	"github.com/gorilla/sessions"
)

var playersRepository repositories.Players = repositories.NewPlayersMemoryRepository()

// PLAYERS

func getWebPlayerId(request *http.Request) int {
	clientSession := request.Context().Value("clientSession").(*sessions.Session)
	wrappedInt, _ := clientSession.Values["clientId"]
	return wrappedInt.(int)
}

func ensurePlayerHasId(request *http.Request, player *repositories.PersistentPlayer) {
	if player.Id == nil {
		id := getWebPlayerId(request)
		player.Id = &id
	}
}

func GetPlayers(response http.ResponseWriter, request *http.Request) {
	players, err := playersRepository.GetPlayers()
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
	player, err := playersRepository.GetPlayerById(id)
	if err != nil {
		if err == repositories.EntityNotExistsErr {
			player = &repositories.PersistentPlayer{
				Player: model.Player{
					Name: "",
				},
				Id: &id,
			}
			player, err = playersRepository.CreatePlayer(*player)
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
	player, err := playersRepository.GetPlayerById(id)
	if err != nil {
		fmt.Printf("error while retrieving player(id=%d): '%v'\n", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, player)
}

func GetPlayersByGame(response http.ResponseWriter, request *http.Request) {
	WriteJsonResponse(response, http.StatusOK, "not implemeted yet")
}

func CreatePlayer(response http.ResponseWriter, request *http.Request) {
	var player repositories.PersistentPlayer
	err := parseJsonFromReader(request.Body, &player)
	if err != nil {
		fmt.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("ParseJsonFromReader(request.Body, &player) = %v %v\n", player, err)
	ensurePlayerHasId(request, &player)
	fmt.Printf("ensurePlayerHasId(request, &player) => %v\n", player)

	created, err := playersRepository.CreatePlayer(player)
	if err != nil {
		fmt.Printf("error while creating Player: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, created)
}

func UpdatePlayer(response http.ResponseWriter, request *http.Request) {
	var player repositories.PersistentPlayer
	err := parseJsonFromReader(request.Body, &player)
	if err != nil {
		fmt.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("ParseJsonFromReader(request.Body, &player) = %v %v\n", player, err)
	ensurePlayerHasId(request, &player)
	fmt.Printf("ensurePlayerHasId(request, &player) => %v\n", player)
	updated, err := playersRepository.UpdatePlayer(player)
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
	err = playersRepository.DeletePlayer(id)
	if err != nil {
		fmt.Printf("error while deleting player(id=%d): '%v'", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
}
