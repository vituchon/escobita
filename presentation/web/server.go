package web

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/gorilla/mux"
)

func Run() {
	portalRouter := buildPortalRouter()
	portalServer := &http.Server{
		Addr:         ":9090",
		Handler:      portalRouter,
		ReadTimeout:  40 * time.Second,
		WriteTimeout: 300 * time.Second,
	}
	fmt.Printf("baby-portal listening on '0.0.0.0%s'", ":9090")
	portalServer.ListenAndServe()

	// TODO (for greater good) : Perhaps we are now in condition to add https://github.com/gorilla/mux#graceful-shutdown
}

func buildPortalRouter() *mux.Router {
	root := mux.NewRouter()
	fileServer := http.FileServer(http.Dir("./"))
	root.PathPrefix("/presentation/web/assets").Handler(fileServer)
	root.NotFoundHandler = http.HandlerFunc(NoMatchingHandler)
	//root.Use(SslRedirect, PortalAccessLogMiddleware, PortalOrgAwareMiddleware)

	Get := BuildSetHandleFunc(root, "GET")
	//Post := BuildSetHandleFunc(root, "POST")
	Get("/healthcheck", healthcheck)

	/*Post("/api/v1/login", controllers.PortalLogin)
	ServePortalHomeAuth := PortalAuthMiddlewareForHome(http.HandlerFunc(controllers.ServePortalHome)).(http.HandlerFunc)
	Get("/home", ServePortalHomeAuth)*/

	/*api := root.PathPrefix("/api/v1").Subrouter()
	api.Use(PortalAuthMiddleware)
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

	http.Redirect(response, request, "/presentation/web/assest/170.png", http.StatusSeeOther)
}
