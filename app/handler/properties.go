package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joaopmgd/github-tag-api/app/model"
	"github.com/joaopmgd/github-tag-api/config"
	"github.com/joaopmgd/github-tag-api/database"
)

// GetAllStarredRepos will recover all the repos starred by an user
func GetAllStarredRepos(config *config.Config, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Validate URL
	URL, err := config.GetStarredReposURL(vars)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	userStarredRepos, err := getUserStarredReposOr404(URL)
	if err != nil {
		config.Log.UnableToRequest(err.Error())
		respondError(w, http.StatusNotFound, "User not found")
		return
	}
	// Recover data from database
	tags := config.DB.GetAllRepoTagsMap(vars["user"])
	respondJSON(w, http.StatusOK, paginate(config, r, createMessageStarredReposSelectedTag(userStarredRepos, tags, r.FormValue("tag"))))
}

func createMessageStarredReposSelectedTag(repos []model.StarredRepoRequest, tags map[int64][]string, selectedTag string) []model.StarredRepoTags {
	starredRepos := []model.StarredRepoTags{}
	for _, repo := range repos {
		if selectedTag == "" || repoHasTag(tags[repo.ID], selectedTag) {
			starredRepos = append(starredRepos, model.StarredRepoTags{
				ID:          repo.ID,
				Name:        repo.Name,
				Description: repo.Description,
				URL:         repo.URL,
				Language:    repo.Language,
				Tags:        tags[repo.ID],
			})
		}
	}
	return starredRepos
}

func repoHasTag(tags []string, selectedTag string) bool {
	for _, tag := range tags {
		if strings.Contains(tag, selectedTag) {
			return true
		}
	}
	return false
}

// getUserStarredReposOr404 gets all user starred repos, or respond the 404 error otherwise
func getUserStarredReposOr404(URL string) ([]model.StarredRepoRequest, error) {
	var userStarredRepos []model.StarredRepoRequest
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

// PostTagStarredRepo post a new tag for a repo
func PostTagStarredRepo(config *config.Config, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Validate body
	var tagData model.TagRequestUpdate
	err := json.NewDecoder(r.Body).Decode(&tagData)
	if err != nil {
		config.Log.CouldNotParseRequestBody(err.Error())
		respondError(w, http.StatusNotFound, "Body must have a JSON key named 'tag' and its value")
		return
	}

	// Validate URL
	URL, err := config.GetStarredReposURL(vars)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate request to github
	userStarredRepos, err := getUserStarredReposOr404(URL)
	if err != nil {
		config.Log.UnableToRequest(err.Error())
		respondError(w, http.StatusNotFound, "User not found")
		return
	}

	// Validate request if id exists
	var repo model.StarredRepoRequest
	for _, starred := range userStarredRepos {
		if vars["repo"] == strconv.FormatInt(starred.ID, 10) {
			repo = starred
		}
	}
	if repo.ID == 0 {
		config.Log.RepoNotFound(vars["repo"])
		respondError(w, http.StatusNotFound, "Repository not found "+vars["repo"])
		return
	}

	// Recover data from database
	tags := config.DB.GetAllRepoTagsByRepoID(vars["user"], repo.ID)

	// If tag already exists return bad request
	for _, tag := range tags {
		if tag.TagName == tagData.TagName {
			config.Log.RepoAlreadyTagged(vars["repo"])
			respondError(w, http.StatusBadRequest, "Repository already has the tag : "+tagData.TagName)
			return
		}
	}
	// Add to database
	config.DB.InsertRepoTagsValue(database.RepoTag{UserID: vars["user"], RepoID: repo.ID, TagName: tagData.TagName})
	config.DB.InsertLanguageTagsValue(database.LanguageTag{Language: repo.Language, TagName: tagData.TagName})
	respondJSON(w, http.StatusOK, model.ResponseOK{Message: "Tag added"})
}

// GetARepoRecommendation will get the most used tags based on a language
func GetARepoRecommendation(config *config.Config, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Validate URL
	URL, err := config.GetStarredReposURL(vars)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	userStarredRepos, err := getUserStarredReposOr404(URL)
	if err != nil {
		config.Log.UnableToRequest(err.Error())
		respondError(w, http.StatusNotFound, "User not found")
		return
	}

	// Validate request if id exists
	var repo model.StarredRepoRequest
	for _, starred := range userStarredRepos {
		if vars["repo"] == strconv.FormatInt(starred.ID, 10) {
			repo = starred
		}
	}
	if repo.ID == 0 {
		config.Log.RepoNotFound(vars["repo"])
		respondError(w, http.StatusNotFound, "Repository not found "+vars["repo"])
		return
	}
	// Recover data from database
	tags := config.DB.GetRecommendationTagByLanguage(repo.Language)
	respondJSON(w, http.StatusOK, model.RecommendedTags{Recommended: addLanguage(repo.Language, tags)})
}

func addLanguage(language string, tags []string) []string {
	for _, tag := range tags {
		if tag == language {
			return tags
		}
	}
	return append(tags, language)
}

// Paginate just picksup a slice from the Response, showing just the page Requested
func paginate(config *config.Config, r *http.Request, starredRepos []model.StarredRepoTags) model.StarredRepoTagsResponse {
	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		limit = 10
	}
	if offset*limit > len(starredRepos)-1 {
		config.Log.PageIsBiggerThanRequestValues(strconv.Itoa(limit), strconv.Itoa(offset))
		return model.StarredRepoTagsResponse{
			StarredRepos:         []model.StarredRepoTags{},
			PageNumber:           offset,
			PageSize:             limit,
			PropertiesTotalCount: len(starredRepos)}
	}
	if offset+limit > len(starredRepos) {
		limit = len(starredRepos)
	}
	selectedRepos := (starredRepos)[offset*limit : (offset*limit)+limit]
	starredReposResponseResponse := model.StarredRepoTagsResponse{
		StarredRepos:         selectedRepos,
		PageNumber:           offset,
		PageSize:             limit,
		PropertiesTotalCount: len(starredRepos),
	}
	return starredReposResponseResponse
}

// HealthStatus checks github and database connectivity
func HealthStatus(config *config.Config, w http.ResponseWriter, r *http.Request) {

	overallStatus := "up"
	databaseStatus := model.DatabaseStatus{Status: "up"}
	err := config.DB.Ping()
	if err != nil {
		overallStatus = "down"
		databaseStatus = model.DatabaseStatus{Status: "down"}
	}
	URL := config.GetHealthStatusURL()
	githubHealth, err := getHealthStatusOr404(URL)
	if err != nil || githubHealth.Status.Indicator != "none" {
		overallStatus = "down"
	}
	respondJSON(w, http.StatusOK, model.AppHealthStatus{Status: overallStatus, Database: databaseStatus, GithubStatus: githubHealth.Status})
}

// getHealthStatusOr404 gets the health status from github
func getHealthStatusOr404(URL string) (model.GithubHealthStatus, error) {
	var health model.GithubHealthStatus
	if err := requestData(&health, URL); err != nil {
		return model.GithubHealthStatus{}, err
	}
	return health, nil
}

// DeleteTagStarredRepo delete a tag for some repo
func DeleteTagStarredRepo(config *config.Config, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Validate body
	var tagData model.TagRequestUpdate
	err := json.NewDecoder(r.Body).Decode(&tagData)
	if err != nil {
		config.Log.CouldNotParseRequestBody(err.Error())
		respondError(w, http.StatusNotFound, "Body must have a JSON key named 'tag' and its value")
		return
	}

	// Validate URL
	URL, err := config.GetStarredReposURL(vars)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate request to github
	userStarredRepos, err := getUserStarredReposOr404(URL)
	if err != nil {
		config.Log.UnableToRequest(err.Error())
		respondError(w, http.StatusNotFound, "User not found")
		return
	}

	// Validate request if id exists
	var repo model.StarredRepoRequest
	for _, starred := range userStarredRepos {
		if vars["repo"] == strconv.FormatInt(starred.ID, 10) {
			repo = starred
		}
	}
	if repo.ID == 0 {
		config.Log.RepoNotFound(vars["repo"])
		respondError(w, http.StatusNotFound, "Repository not found "+vars["repo"])
		return
	}

	config.DB.DeleteRepoTagsValue(database.RepoTag{UserID: vars["user"], RepoID: repo.ID, TagName: tagData.TagName})
	respondJSON(w, http.StatusOK, model.ResponseOK{Message: "Tag Deleted"})
}
