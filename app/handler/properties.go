package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/github-tag-api/app/model"
	"github.com/github-tag-api/config"
	"github.com/gorilla/mux"
)

// GetAllStarredRepos will recover all the repos starred by an user
func GetAllStarredRepos(config *config.Config, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["user"] == "" {
		config.Log.NoUserSet(vars["user"])
		respondError(w, http.StatusBadRequest, "No user passed with the response")
		return
	}
	URL, err := config.GetStarredReposURL(vars)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	userStarredRepos, err := getUserStarredReposOr404(URL)
	if err != nil {
		config.Log.UnableToRequest(err.Error())
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, userStarredRepos)
}

// getUserStarredReposOr404 gets all user starred repos, or respond the 404 error otherwise
func getUserStarredReposOr404(URL string) ([]model.StarredRepo, error) {
	var userStarredRepos []model.StarredRepo
	if err := requestData(&userStarredRepos, URL); err != nil {
		return nil, err
	}
	return userStarredRepos, nil
}

func requestData(target interface{}, URL string) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(URL)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
