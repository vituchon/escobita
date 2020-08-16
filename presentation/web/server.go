package web

// trying to use one file for all kinds controllers and content delivery related code

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"time"

	"local/escobita/presentation/util"
	"local/escobita/presentation/web/controllers"
	"local/escobita/presentation/web/services"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
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

var ClientSessions *sessions.CookieStore

func retrieveCookieStoreKey() (key []byte, err error) {
	if util.FileExists(storeKeyFilePath) {
		key, err = ioutil.ReadFile(storeKeyFilePath)
		fmt.Printf("Using existing key %s\n", string(key))
	} else {
		key = securecookie.GenerateRandomKey(32)
		ioutil.WriteFile(storeKeyFilePath, key, 0644)
		fmt.Printf("Generated new key %s and stored at %s", string(key), storeKeyFilePath)
	}
	return
}

func StartWebServer() {
	key, err := retrieveCookieStoreKey()
	if err != nil {
		fmt.Printf("Unexpected error while retrieving cookie store key: %v", err)
		return
	}
	ClientSessions = sessions.NewCookieStore(key)

	router := buildRouter()
	server := &http.Server{
		Addr:         ":9090",
		Handler:      router,
		ReadTimeout:  40 * time.Second,
		WriteTimeout: 300 * time.Second,
	}
	fmt.Printf("escobita web server listening at port %v", server.Addr)
	server.ListenAndServe()

	// TODO (for greater good) : Perhaps we are now in condition to add https://github.com/gorilla/mux#graceful-shutdown
}

func buildRouter() *mux.Router {
	root := mux.NewRouter()
	fileServer := http.FileServer(http.Dir("./"))
	root.PathPrefix("/presentation/web/assets").Handler(fileServer)
	root.NotFoundHandler = http.HandlerFunc(NoMatchingHandler)
	//root.Use(SslRedirect, AccessLogMiddleware, OrgAwareMiddleware)
	root.Use(ClientSessionAwareMiddleware, AccessLogMiddleware)

	Get := BuildSetHandleFunc(root, "GET")
	//Post := BuildSetHandleFunc(root, "POST")
	Get("/", serveRoot)
	Get("/healthcheck", controllers.Healthcheck)

	/*Post("/api/v1/login", controllers.Login)
	ServeHomeAuth := AuthMiddlewareForHome(http.HandlerFunc(controllers.ServeHome)).(http.HandlerFunc)
	Get("/home", ServeHomeAuth)*/

	api := root.PathPrefix("/api/v1").Subrouter()
	//api.Use(AuthMiddleware)
	apiGet := BuildSetHandleFunc(api, "GET")
	apiPost := BuildSetHandleFunc(api, "POST")
	apiPut := BuildSetHandleFunc(api, "PUT")
	apiDelete := BuildSetHandleFunc(api, "DELETE")

	apiGet("/games", controllers.GetGames)
	apiGet("/games/{id:[0-9]+}", controllers.GetGameById)
	apiPost("/games", controllers.CreateGame)
	apiPut("/games/{id:[0-9]+}", controllers.UpdateGame)
	apiDelete("/games/{id:[0-9]+}", controllers.UpdateGame)

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

	if request.URL.Path == "/favicon.ico" { // avoids to trigger another request to landing or login on the "silent" http request by chrome to get an icon! I guess i could tell chrome for ubuntu that redirection for an icon can create more and bigger troubles than solutions... i mean nobody dies for an icon... for now...
		response.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(response, request, "/presentation/web/assets/images/logo.png", http.StatusSeeOther)
}

// Adds a logging handler for logging each request's in Apache Common Log Format (CLF).
// With this middleware we ensure that each requests will be, at least, logged once.
func AccessLogMiddleware(h http.Handler) http.Handler {
	loggingHandler := handlers.LoggingHandler(services.Logger.Out, h)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggingHandler.ServeHTTP(w, r)
	})
}

func ClientSessionAwareMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		clientSession := getOrCreateClientSession(request)
		err := saveClientSession(request, response, clientSession)
		if err != nil {
			fmt.Printf("error while saving client session: %v", err)
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.ServeHTTP(response, request)
	})
}

func getOrCreateClientSession(request *http.Request) *sessions.Session {
	clientSession, err := ClientSessions.Get(request, "client_session")
	if err != nil {
		fmt.Printf("error while retrieving 'client_session' from session store: %+v\n", err)
	}
	if clientSession.IsNew {
		fmt.Print("creating new session\n")
	} else {
		fmt.Print("using existing session\n")
	}
	return clientSession
}

func saveClientSession(request *http.Request, response http.ResponseWriter, clientSession *sessions.Session) error {
	return ClientSessions.Save(request, response, clientSession)
}

// Dev notes: the request context has the organization due to the ContextAwareMiddle, so there will be always a valid portal's client session when invoking this function
func serveRoot(response http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("presentation/web/assets/html/root.html")
	if err != nil {
		fmt.Printf("Error while parsing template : %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	t.Execute(response, nil)
}
