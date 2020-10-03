package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
