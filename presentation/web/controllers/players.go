package controllers

import (
	"log"
	"net/http"
	"sync"

	"github.com/vituchon/escobita/model"

	"github.com/vituchon/escobita/presentation/web/services"
	"github.com/vituchon/escobita/repositories"
)

var playersRepository repositories.Players = repositories.NewPlayersMemoryRepository()

// PLAYERS

func GetPlayers(response http.ResponseWriter, request *http.Request) {
	players, err := playersRepository.GetPlayers()
	if err != nil {
		log.Printf("error while retrieving players : '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, players)
}

var getClientPlayerMutex sync.Mutex

// Gets the web client's correspondant player (There is only ONE player per client!)
func GetClientPlayer(response http.ResponseWriter, request *http.Request) {
	getClientPlayerMutex.Lock()
	defer getClientPlayerMutex.Unlock()
	id := services.GetClientId(request)
	player, err := playersRepository.GetPlayerById(id)
	if err != nil {
		if err == repositories.EntityNotExistsErr {
			name := ""
			paramName, err := ParseSingleStringUrlQueryParam(request, "name")
			if err != nil && paramName != nil {
				name = *paramName
			}
			player, err = createPlayer(id, name) // create player in memory
			if err != nil {
				log.Printf("error while getting client player : '%v'", err)
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, err = playersRepository.CreatePlayer(*player) // saves player in a persistent storage
			if err != nil {
				log.Printf("error while getting client player : '%v'", err)
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Printf("Creating new player %+v for ip=%s\n", player, request.RemoteAddr)
		} else {
			log.Printf("error while getting client player : '%v'", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("Using existing player %+v for ip=%s \n", player, request.RemoteAddr)
	}
	WriteJsonResponse(response, http.StatusOK, player)
}

func createPlayer(id int, name string) (*repositories.PersistentPlayer, error) {
	err := model.ValidateName(name)
	if err != nil {
		return nil, err
	}

	return &repositories.PersistentPlayer{
		Name: name,
		Id:   id,
	}, nil
}

func GetPlayerById(response http.ResponseWriter, request *http.Request) {
	id, err := ParseRouteParamAsInt(request, "id")
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	player, err := playersRepository.GetPlayerById(id)
	if err != nil {
		log.Printf("error while retrieving player(id=%d): '%v'\n", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, player)
}

func GetPlayersByGame(response http.ResponseWriter, request *http.Request) {
	WriteJsonResponse(response, http.StatusBadRequest, "Endpoint not implemeted")
}

func UpdatePlayer(response http.ResponseWriter, request *http.Request) {
	var player repositories.PersistentPlayer
	err := parseJsonFromReader(request.Body, &player)
	if err != nil {
		log.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	updated, err := playersRepository.UpdatePlayer(player)
	if err != nil {
		log.Printf("error while updating Player: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, updated)
}

func DeletePlayer(response http.ResponseWriter, request *http.Request) {
	id, err := ParseRouteParamAsInt(request, "id")
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	err = playersRepository.DeletePlayer(id)
	if err != nil {
		log.Printf("error while deleting player(id=%d): '%v'", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
}
