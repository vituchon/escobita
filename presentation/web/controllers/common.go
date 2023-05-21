package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Healthcheck(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
}

const ServerVersion = "0.0.1"

func Version(response http.ResponseWriter, request *http.Request) {
	/*response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")*/
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

func ParseRouteParamAsInt(request *http.Request, name string) (int, error) {
	rawValue := mux.Vars(request)[name]
	intValue, err := strconv.Atoi(rawValue)
	if err != nil {
		errMsg := fmt.Sprintf("Can not parse route param as integer from '%s'", intValue)
		return 0, errors.New(errMsg)
	}
	return intValue, nil
}

var (
	UrlParamNotFoundErr = errors.New("No url param present with the given name")
)

// Gets an integer url's query param with the given name
func UrlQueryIntParam(request *http.Request, name string) (*int, error) {
	param := request.URL.Query().Get(name)
	if len(param) == 0 {
		return nil, UrlParamNotFoundErr
	}
	value, err := strconv.Atoi(param)
	return &value, err
}
