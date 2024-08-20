package controllers

import (
	"log"
	"net/http"

	"github.com/vituchon/escobita/presentation/web/services"
)

func GetMessages(response http.ResponseWriter, request *http.Request) {
	messages, err := services.GetMessages()
	if err != nil {
		log.Printf("error while retrieving messages : '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, messages)
}

func GetMessageById(response http.ResponseWriter, request *http.Request) {
	id, err := ParseRouteParamAsInt(request, "id")
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	message, err := services.GetMessageById(id)
	if err != nil {
		log.Printf("error while retrieving message(id=%d): '%v'", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, message)
}

func GetMessagesByGame(response http.ResponseWriter, request *http.Request) {
	id, err := ParseRouteParamAsInt(request, "id")
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	since, err := ParseSingleIntegerUrlQueryParam(request, "since")
	if err != nil {
		if err == UrlQueryParamNotFoundErr {
			var zero = 0
			since = &zero
		} else {
			log.Printf("error parsing url param 'since': '%v'", err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	messages, err := services.GetMessagesByGameAndTime(id, int64(*since))
	if err != nil {
		log.Printf("error while retrieving messages for game(id=%d): '%v'", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, messages)
}

func CreateMessage(response http.ResponseWriter, request *http.Request) {
	var message services.PersistentWebMessage
	err := parseJsonFromReader(request.Body, &message)
	if err != nil {
		log.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	created, err := services.CreateMessage(message)
	if err != nil {
		log.Printf("error while creating Message: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, created)
}

func UpdateMessage(response http.ResponseWriter, request *http.Request) {
	var message services.PersistentWebMessage
	err := parseJsonFromReader(request.Body, &message)
	if err != nil {
		log.Printf("error reading request body: '%v'", err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	updated, err := services.UpdateMessage(message)
	if err != nil {
		log.Printf("error while updating Message: '%v'", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	WriteJsonResponse(response, http.StatusOK, updated)
}

func DeleteMessage(response http.ResponseWriter, request *http.Request) {
	id, err := ParseRouteParamAsInt(request, "id")
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	err = services.DeleteMessage(id)
	if err != nil {
		log.Printf("error while deleting message(id=%d): '%v'", id, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
}
