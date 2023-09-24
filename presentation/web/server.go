package web

// The presentation layer contains all resources concerned with creating an application interface
// Contains code designed to be used for http rest based api interface

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/vituchon/escobita/presentation/util"
	"github.com/vituchon/escobita/presentation/web/controllers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	//	embed "embed"
)

// TODO : nice implement something like this
/*
type Server struct {
	store *CookieStore
}

func (s Server) Startup() {

}

func (s Server) Shutdown() {

}
*/

const (
	storeKeyFilePath = ".ss" // the file were the actual key is stored
)

func retrieveCookieStoreKey(filepath string) (key []byte, err error) {
	if util.FileExists(filepath) {
		key, err = ioutil.ReadFile(storeKeyFilePath)
		//fmt.Printf("Using existing key %s\n", string(key))
	} else {
		key = securecookie.GenerateRandomKey(32)
		ioutil.WriteFile(storeKeyFilePath, key, 0644)
		//fmt.Printf("Generated new key %s and stored at %s\n", string(key), storeKeyFilePath)
	}
	return
}

/*
//go:embed assets/*
var assets embed.FS*/

func StartServer() {
	//file, err := assets.Open("assets/html/root.html")
	//bytes, err := ioutil.ReadFile("./pepe.txt")

	//fmt.Println(err, string(bytes))
	//fileBytes, err := ioutil.ReadAll(file)

	key, err := retrieveCookieStoreKey(storeKeyFilePath)
	if err != nil {
		fmt.Printf("Unexpected error while retrieving cookie store key: %v", err)
		return
	}
	controllers.NewSessionStore(key)

	router := buildRouter()
	server := &http.Server{
		Addr:         ":9090",
		Handler:      router,
		ReadTimeout:  40 * time.Second,
		WriteTimeout: 300 * time.Second,
	}
	fmt.Printf("escobita web server listening at port %v", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("Unexpected error initiliazing web server: ", err)
	}

	// TODO (for greater good) : Perhaps we are now in condition to add https://github.com/gorilla/mux#graceful-shutdown
}

func buildRouter() *mux.Router {
	root := mux.NewRouter()
	// TODO : word "presentation" in the path may be redudant, perpahs using just "assets" would be enought!
	// BEFORE go:embed
	fileServer := http.FileServer(http.Dir("./"))
	root.PathPrefix("/presentation/web/assets").Handler(fileServer)

	// AFTER go:embed
	/*fileServer := http.FileServer(http.FS(assets))
	root.PathPrefix("/presentation/web/").Handler(http.StripPrefix("/presentation/web/", fileServer))*/

	root.NotFoundHandler = http.HandlerFunc(NoMatchingHandler)
	//root.Use(SslRedirect, AccessLogMiddleware, OrgAwareMiddleware)
	root.Use(ClientSessionAwareMiddleware)

	Get := BuildSetHandleFunc(root, "GET")
	//Post := BuildSetHandleFunc(root, "POST")
	Get("/", serveRoot)
	Get("/healthcheck", controllers.Healthcheck)
	Get("/version", controllers.Version)

	Get("/adquire-ws", controllers.AdquireWebSocket)
	Get("/release-ws", controllers.ReleaseWebSocket)
	Get("/debug-ws", controllers.DebugWebSockets)
	Get("/send-message-ws", controllers.SendMessageWebSocket)

	Get("/send-message-all-ws", controllers.SendMessageAllWebSockets)
	Get("/release-broken-ws", controllers.ReleaseBrokenWebSockets)
	Get("/release-all-ws", controllers.ReleaseAllWebSockets)

	/*Post("/api/v1/login", controllers.Login)
	ServeHomeAuth := AuthMiddlewareForHome(http.HandlerFunc(controllers.ServeHome)).(http.HandlerFunc)
	Get("/home", ServeHomeAuth)*/

	api := root.PathPrefix("/api/v1").Subrouter()
	api.Use(AccessLogMiddleware) // only logs api calls
	//api.Use(AuthMiddleware)
	apiGet := BuildSetHandleFunc(api, "GET")
	apiPost := BuildSetHandleFunc(api, "POST")
	apiPut := BuildSetHandleFunc(api, "PUT")
	apiDelete := BuildSetHandleFunc(api, "DELETE")

	apiGet("/games", controllers.GetGames)
	apiGet("/games/{id:[0-9]+}", controllers.GetGameById)
	apiPost("/games", controllers.CreateGame)
	apiPost("/games/{id:[0-9]+}/message", controllers.SendMessage)
	apiPut("/games/{id:[0-9]+}", controllers.UpdateGame)
	apiDelete("/games/{id:[0-9]+}", controllers.DeleteGame)
	apiPost("/games/{id:[0-9]+}/resume", controllers.ResumeGame)
	apiPost("/games/{id:[0-9]+}/perform-take-action", controllers.PerformTakeAction)
	apiPost("/games/{id:[0-9]+}/perform-drop-action", controllers.PerformDropAction)
	apiGet("/games/{id:[0-9]+}/calculate-stats", controllers.CalculateGameStats)

	apiGet("/games/{id:[0-9]+}/bind-ws", controllers.BindClientWebSocketToGame)
	apiGet("/games/{id:[0-9]+}/unbind-ws", controllers.UnbindClientWebSocketInGame)

	apiGet("/players", controllers.GetPlayers)
	apiGet("/players-by-game", controllers.GetPlayersByGame) // TODO: Implement if necessary...
	apiGet("/player", controllers.GetClientPlayer)
	apiGet("/players/{id:[0-9]+}", controllers.GetPlayerById)
	apiPost("/players", controllers.CreatePlayer)
	apiPut("/players/{id:[0-9]+}", controllers.UpdatePlayer)

	apiGet("/messages", controllers.GetMessages) // TODO : add optional parameter "since", default beign server start up time
	apiGet("/messages/{id:[0-9]+}", controllers.GetMessageById)
	apiGet("/messages/get-by-game/{id:[0-9]+}", controllers.GetMessagesByGame)
	apiPost("/messages", controllers.CreateMessage)
	apiPut("/messages/{id:[0-9]+}", controllers.UpdateMessage)
	apiDelete("/messages/{id:[0-9]+}", controllers.DeleteMessage)

	return root
}

type setHandlerFunc func(path string, f http.HandlerFunc)

// Creates a function for register a handler for a path for the given router and http methods
func BuildSetHandleFunc(router *mux.Router, methods ...string) setHandlerFunc {
	return func(path string, f http.HandlerFunc) {
		router.HandleFunc(path, f).Methods(methods...)
	}
}

func NoMatchingHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("No maching route for " + request.URL.Path)
	response.WriteHeader(http.StatusNotFound)

	/*if request.URL.Path == "/favicon.ico" { // avoids to trigger another request to landing or login on the "silent" http request by chrome to get an icon! I guess i could tell chrome for ubuntu that redirection for an icon can create more and bigger troubles than solutions... i mean nobody dies for an icon... for now...
		response.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(response, request, "/presentation/web/assets/images/logo.png", http.StatusSeeOther)*/
}

// Adds a logging handler for logging each request's in Apache Common Log Format (CLF).
// With this middleware we ensure that each requests will be, at least, logged once.
func AccessLogMiddleware(h http.Handler) http.Handler {
	loggingHandler := handlers.LoggingHandler(os.Stdout, h)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggingHandler.ServeHTTP(w, r)
	})
}

func ClientSessionAwareMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		clientSession, err := controllers.GetOrCreateClientSession(request)
		if err != nil {
			fmt.Printf("error while getting client session: %v", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = controllers.SaveClientSession(request, response, clientSession)
		if err != nil {
			fmt.Printf("error while saving client session: %v", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(request.Context(), "clientSession", clientSession)
		h.ServeHTTP(response, request.WithContext(ctx))
	})
}

// Dev notes: the request context has the organization due to the ContextAwareMiddle, so there will be always a valid portal's client session when invoking this function
func serveRoot(response http.ResponseWriter, request *http.Request) {
	//t, err := template.ParseFS(assets, "assets/html/root.html")
	t, err := template.ParseFiles("presentation/web/assets/html/root.html")
	if err != nil {
		fmt.Printf("Error while parsing template : %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	t.Execute(response, nil)
}
