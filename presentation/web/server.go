package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"time"

	"local/escobita/presentation/util"

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
	root.Use(ClientSessionAwareMiddleware)

	Get := BuildSetHandleFunc(root, "GET")
	//Post := BuildSetHandleFunc(root, "POST")
	Get("/healthcheck", healthcheck)

	/*Post("/api/v1/login", controllers.Login)
	ServeHomeAuth := AuthMiddlewareForHome(http.HandlerFunc(controllers.ServeHome)).(http.HandlerFunc)
	Get("/home", ServeHomeAuth)*/

	/*api := root.PathPrefix("/api/v1").Subrouter()
	api.Use(AuthMiddleware)
	apiGet := BuildSetHandleFunc(api, "GET")
	apiPost := BuildSetHandleFunc(api, "POST")
	apiPut := BuildSetHandleFunc(api, "PUT")
	apiDelete := BuildSetHandleFunc(api, "DELETE")*/

	return root
}

func healthcheck(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
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

	if request.URL.Path == "/favicon.ico" { // avoids to trigger another request to landing or login on the "silent" http request by chrome to get an icon! I guess i could tell chrome for ubuntu that redirection for an icon can create more and bigger troubles than solutions... i mean nobody dies for an icon... for now...
		response.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(response, request, "/presentation/web/assets/170.png", http.StatusSeeOther)
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
