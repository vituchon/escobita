package controllers

import (
	"fmt"
	"local/escobita/model"
	"local/escobita/presentation/web/services"
	"net/http"
	"strconv"
)

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
	game.PlayerId = getWebPlayerId(request) // asign owner
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

func ResumeGame(response http.ResponseWriter, request *http.Request) {
	var game services.WebGame
	err := ParseJsonFromReader(request.Body, &game)
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
