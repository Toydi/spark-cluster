package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(handler *APIHandler) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	var publicRoutes = Routes{
		{
			"Authorize",
			"GET",
			"/oauth/callback",
			handler.Authorize,
		},
		{
			"CreateTerminal",
			"GET",
			"/apis/terminal",
			handler.CreateTerminal,
		},
		{
			"UploadFile",
			"GET",
			"/apis/file",
			handler.UploadFile,
		},
		{
			"GetLog",
			"GET",
			"/apis/log",
			handler.GetLog,
		},
	}

	var privateRoutes = Routes{
		{
			"ListDatasets",
			"GET",
			"/apis/dataset",
			handler.ListDatasets,
		},
		{
			"CreateDataset",
			"POST",
			"/apis/dataset",
			handler.CreateDataset,
		},		
		{
			"ListSparkCluster",
			"GET",
			"/apis/sparkcluster",
			handler.ListSparkCluster,
		},
		{
			"CreateSparkCluster",
			"POST",
			"/apis/sparkcluster",
			handler.CreateSparkCluster,
		},
		{
			"DeleteSparkCluster",
			"DELETE",
			"/apis/sparkcluster/{sparkcluster}",
			handler.DeleteSparkCluster,
		},
		{
			"UpdateSparkCluster",
			"PUT",
			"/apis/sparkcluster/{sparkcluster}",
			handler.UpdateSparkCluster,
		},		
		{
			"CurrentUser",
			"GET",
			"/apis/user",
			handler.CurrentUser,
		},
	}
	// The public route is always accessible
	for _, route := range publicRoutes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	// The private route is only accessible if the user has a valid access_token.
	// We are chaining the middleware into the negroni handler function which will
	// check for a valid token.
	for _, route := range privateRoutes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(negroni.New(
				negroni.HandlerFunc(handler.AuthMiddleware),
				negroni.Wrap(route.HandlerFunc)))
	}

	// Handle websocket routes with path prefix router
	// WebSocket need multi routes for a service.
	router.PathPrefix("/terminal/ws").Handler(NewTerminal(handler.kubeClient, handler.kubeConfig))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization"},
	})

	return c.Handler(router)
}
