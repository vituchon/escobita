package web

// The presentation layer contains all resources concerned with creating an application interface
// Contains code designed to be used for http rest based api interface

import (
	"context"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/vituchon/escobita/presentation/util"
	"github.com/vituchon/escobita/presentation/web/controllers"

	//embed "embed"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
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
		//log.Printf("Using existing key %s\n", string(key))
	} else {
		key = securecookie.GenerateRandomKey(32)
		ioutil.WriteFile(storeKeyFilePath, key, 0644)
		//log.Printf("Generated new key %s and stored at %s\n", string(key), storeKeyFilePath)
	}
	return
}

/*
//go:embed assets/*
var assets embed.FS*/

func StartServer() {
	//file, err := assets.Open("assets/html/root.html")
	//bytes, err := ioutil.ReadFile("./pepe.txt")

	//log.Println(err, string(bytes))
	//fileBytes, err := ioutil.ReadAll(file)

	key, err := retrieveCookieStoreKey(storeKeyFilePath)
	if err != nil {
		log.Printf("Unexpected error while retrieving cookie store key: %v", err)
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
	log.Printf("Escobita web server listening at port %v", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Println("Unexpected error initiliazing escobita web server: ", err)
	}

	// TODO (for greater good) : Perhaps we are now in condition to add https://github.com/gorilla/mux#graceful-shutdown
}

func buildRouter() *mux.Router {
	router := mux.NewRouter()
	//router.Use(AccessLogMiddleware) // no logging at this level as assets' requests would be logged thus creating a lot of unnecessary log messages
	router.NotFoundHandler = http.HandlerFunc(NoMatchingHandler)

	// BEFORE go:embed
	assetsFileServer := http.FileServer(http.Dir("./"))
	assetsRouter := router.PathPrefix("/presentation/web/assets").Subrouter()
	assetsRouter.PathPrefix("/").Handler(assetsFileServer)

	// AFTER go:embed
	/*assetsFileServer := http.FileServer(http.FS(assets))
	assetsRouter := router.PathPrefix("/presentation/web/assets").Subrouter()
	assetsRouter.PathPrefix("/").Handler(http.StripPrefix("/presentation/web/", assetsFileServer))*/

	rootRouter := router.PathPrefix("/").Subrouter()
	rootRouter.Use(AccessLogMiddleware, ClientSessionAwareMiddleware)

	rootGet := BuildSetHandleFunc(rootRouter, "GET")
	rootGet("/", serveRoot)
	rootGet("/healthcheck", controllers.Healthcheck)
	rootGet("/version", controllers.Version)

	rootGet("/adquire-ws", controllers.AdquireWebSocket)
	rootGet("/release-ws", controllers.ReleaseWebSocket)
	rootGet("/debug-ws", controllers.DebugWebSockets)
	rootGet("/send-message-ws", controllers.SendMessageWebSocket)

	rootGet("/send-message-all-ws", controllers.SendMessageAllWebSockets)
	rootGet("/release-broken-ws", controllers.ReleaseBrokenWebSockets)
	rootGet("/release-all-ws", controllers.ReleaseAllWebSockets)

	apiRouter := rootRouter.PathPrefix("/api/v1").Subrouter()
	apiGet := BuildSetHandleFunc(apiRouter, "GET")
	apiPost := BuildSetHandleFunc(apiRouter, "POST")
	apiPut := BuildSetHandleFunc(apiRouter, "PUT")
	apiDelete := BuildSetHandleFunc(apiRouter, "DELETE")

	apiGet("/games", controllers.GetGames)
	apiGet("/games/{id:[0-9]+}", controllers.GetGameById)
	apiPost("/games", controllers.CreateGame)
	apiPost("/games/{id:[0-9]+}/message", controllers.SendMessage)
	//apiPut("/games/{id:[0-9]+}", controllers.UpdateGame)
	apiDelete("/games/{id:[0-9]+}", controllers.DeleteGame)
	apiPost("/games/{id:[0-9]+}/start", controllers.StartGame)
	apiPost("/games/{id:[0-9]+}/join", controllers.JoinGame)
	apiPost("/games/{id:[0-9]+}/quit", controllers.QuitGame)
	apiPost("/games/{id:[0-9]+}/perform-take-action", controllers.PerformTakeAction)
	apiPost("/games/{id:[0-9]+}/perform-drop-action", controllers.PerformDropAction)
	apiGet("/games/{id:[0-9]+}/calculate-stats", controllers.CalculateGameStats)

	apiGet("/games/{id:[0-9]+}/bind-ws", controllers.BindClientWebSocketToGame)
	apiGet("/games/{id:[0-9]+}/unbind-ws", controllers.UnbindClientWebSocketInGame)

	apiGet("/players", controllers.GetPlayers)
	apiGet("/players-by-game", controllers.GetPlayersByGame) // TODO: Implement if necessary...
	apiGet("/player", controllers.GetClientPlayer)
	apiGet("/players/{id:[0-9]+}", controllers.GetPlayerById)
	apiPut("/players/{id:[0-9]+}", controllers.UpdatePlayer)

	apiGet("/messages", controllers.GetMessages) // TODO : add optional parameter "since", default beign server start up time
	apiGet("/messages/{id:[0-9]+}", controllers.GetMessageById)
	apiGet("/messages/get-by-game/{id:[0-9]+}", controllers.GetMessagesByGame)
	apiPost("/messages", controllers.CreateMessage)
	apiPut("/messages/{id:[0-9]+}", controllers.UpdateMessage)
	apiDelete("/messages/{id:[0-9]+}", controllers.DeleteMessage)

	return router
}

type setHandlerFunc func(path string, f http.HandlerFunc)

// Creates a function for register a handler for a path for the given router and http methods
func BuildSetHandleFunc(router *mux.Router, methods ...string) setHandlerFunc {
	return func(path string, f http.HandlerFunc) {
		router.HandleFunc(path, f).Methods(methods...)
	}
}

func NoMatchingHandler(response http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/favicon.ico" { // don't log this
		log.Println("No maching route for " + request.URL.Path)
	}
	response.WriteHeader(http.StatusNotFound)
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
			log.Printf("error while getting client session: %v", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = controllers.SaveClientSession(request, response, clientSession)
		if err != nil {
			log.Printf("error while saving client session: %v", err)
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
		log.Printf("Error while parsing template : %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	t.Execute(response, nil)
}
