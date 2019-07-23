package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaopmgd/github-tag-api/app/handler"
	"github.com/joaopmgd/github-tag-api/config"
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
	a.Post("/repos/{user}/starred/{repo}", a.PostTagStarredRepo)
	a.Delete("/repos/{user}/starred/{repo}", a.DeleteTagStarredRepo)
	a.Get("/repos/{user}/starred/{repo}/recommendation", a.GetARepoRecommendation)
	a.Get("/health", a.HealthStatus)
}

// Get Wrap the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post Wrap the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Delete Wrap the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// GetAllStarredRepos Handlers to manage all starred repos
func (a *App) GetAllStarredRepos(w http.ResponseWriter, r *http.Request) {
	handler.GetAllStarredRepos(a.Config, w, r)
}

// PostTagStarredRepo Handlers to post a new tag to a repo
func (a *App) PostTagStarredRepo(w http.ResponseWriter, r *http.Request) {
	handler.PostTagStarredRepo(a.Config, w, r)
}

// DeleteTagStarredRepo Handlers to post a new tag to a repo
func (a *App) DeleteTagStarredRepo(w http.ResponseWriter, r *http.Request) {
	handler.DeleteTagStarredRepo(a.Config, w, r)
}

// GetARepoRecommendation Handlers to get recommendations based on a language
func (a *App) GetARepoRecommendation(w http.ResponseWriter, r *http.Request) {
	handler.GetARepoRecommendation(a.Config, w, r)
}

// HealthStatus returns the health status of the app
func (a *App) HealthStatus(w http.ResponseWriter, r *http.Request) {
	handler.HealthStatus(a.Config, w, r)
}

// Run the app on it's router
func (a *App) Run(host string) {
	a.Config.Log.ListeningPort(host)
	log.Fatal(http.ListenAndServe(host, a.Router))
}
