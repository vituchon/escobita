package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"local/escobita/presentation/web/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Healthcheck(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
}

func WriteJsonResponse(response http.ResponseWriter, status int, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("error while mashalling object %v, trace: %+v", data, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("Escribiendo estos bytes %v\ncomo string son %v\n", bytes, string(bytes))
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

// game controller

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
	fmt.Println("HEMOS LLEGADO!!")
	var game services.WebGame
	fmt.Printf("request.Body = %v\n", request.Body)
	err := ParseJsonFromReader(request.Body, &game)
	fmt.Printf("Tenemos esto %v %v\n", game, err)
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
	fmt.Printf("JUEGO CREADO ES %+v\n", created)
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
