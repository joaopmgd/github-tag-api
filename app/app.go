package app

import (
	"net/http"

	"github.com/github-tag-api/app/handler"
	"github.com/github-tag-api/config"
	"github.com/gorilla/mux"
)

var log *config.StandardLogger

// App has router
type App struct {
	Config *config.Config
	Router *mux.Router
}

// Initialize with predefined configuration and check for environment variables
func (a *App) Initialize(config *config.Config) {
	a.Config = config
	a.Config.Log.EnvVariablesData()
	a.Config.Log.InitFunction("Github Tag API")
	a.Router = mux.NewRouter()
	a.setRouters()
	a.Router.Use(loggingMiddleware)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.MiddlewareRequest(r)
		next.ServeHTTP(w, r)
	})
}

// Set all required routers
func (a *App) setRouters() {
	log = a.Config.Log
	a.Config.Log.SettingUpRouters()
	a.Get("/repos/{user}/starred", a.GetAllStarredRepos)
}

// Get Wrap the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// GetAllStarredRepos Handlers to manage all starred repos
func (a *App) GetAllStarredRepos(w http.ResponseWriter, r *http.Request) {
	handler.GetAllStarredRepos(a.Config, w, r)
}

// Run the app on it's router
func (a *App) Run(host string) {
	a.Config.Log.ListeningPort(host)
	log.Fatal(http.ListenAndServe(host, a.Router))
}
