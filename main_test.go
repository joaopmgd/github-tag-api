package main_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joaopmgd/github-tag-api/app"
	"github.com/joaopmgd/github-tag-api/config"
)

var a app.App

func init() {
	config := config.GetConfig()
	a.Initialize(config)
}

// TestGetStarredRepos Tests for getting starred repos
func TestGetStarredRepos(t *testing.T) {
	tt := map[string]struct {
		user           string
		tag            string
		responseStatus int
		responseBody   string
	}{
		"correct_user":  {"joaopmgd", "", http.StatusOK, ""},
		"no_user_found": {"joaopmgdjoaopmgdjoaopmgdjoaopmgd", "", http.StatusNotFound, "{\"error\":\"User not found\"}"},
		"no_user":       {"", "", http.StatusNotFound, "{\"error\":\"User not found\"}"},
		"no_tag_found":  {"joaopmgd", "do_not_exist", http.StatusOK, ""},
		"tag_found":     {"joaopmgd", "exist", http.StatusOK, ""},
	}
	for testName, tc := range tt {

		req, err := http.NewRequest("GET", "/repos/{user}/starred", nil)
		if err != nil {
			t.Fatal(err)
		}
		req = mux.SetURLVars(req, map[string]string{
			"user": tc.user,
		})

		q := req.URL.Query()
		q.Add("tag", tc.tag)
		q.Add("limit", "1")

		req.URL.RawQuery = q.Encode()
		response := executeRequestTest(req, a.GetAllStarredRepos)

		if response.Code != tc.responseStatus ||
			(response.Code != http.StatusOK && response.Body.String() != tc.responseBody) {
			t.Errorf("\nHandler %s\nFor user %s\nGot Status %v and Body %s\nWant Status %v and Body %s",
				testName, tc.user, response.Code, response.Body.String(), tc.responseStatus, tc.responseBody)
		}
	}
}

// TestPostTagStarredRepo Tests for addding new tags for repos
func TestPostTagStarredRepo(t *testing.T) {
	tt := map[string]struct {
		user           string
		repo           int64
		requestBody    string
		responseStatus int
		responseBody   string
	}{
		"existing_user_and_repo": {"joaopmgd", 10866521, "{\"tag\": \"test\"}", http.StatusOK, "{\"Message\":\"Tag added\"}"},
		"repeated_tag":           {"joaopmgd", 10866521, "{\"tag\": \"test\"}", http.StatusBadRequest, "{\"error\":\"Repository already has the tag : test\"}"},
		"user_not_found":         {"joaopmgdjoaopmgdjoaopmgdjoaopmgd", 10866521, "{\"tag\": \"test\"}", http.StatusNotFound, "{\"error\":\"User not found\"}"},
		"no_body":                {"joaopmgd", 10866521, "", http.StatusNotFound, "{\"error\":\"Body must have a JSON key named 'tag' and its value\"}"},
		"repo_not_found":         {"joaopmgd", 999999999, "{\"tag\": \"test\"}", http.StatusNotFound, "{\"error\":\"Repository not found 999999999\"}"},
	}
	for testName, tc := range tt {

		req, err := http.NewRequest("POST", "/repos/{user}/starred/{repo}", strings.NewReader(tc.requestBody))
		if err != nil {
			t.Fatal(err)
		}
		req = mux.SetURLVars(req, map[string]string{
			"user": tc.user,
			"repo": strconv.FormatInt(tc.repo, 10),
		})
		response := executeRequestTest(req, a.PostTagStarredRepo)

		if response.Code != tc.responseStatus ||
			response.Body.String() != tc.responseBody {
			t.Errorf("\nHandler %s\nFor user %s\nGot Status %v and Body %s\nWant Status %v and Body %s",
				testName, tc.user, response.Code, response.Body.String(), tc.responseStatus, tc.responseBody)
		}
	}
	removeCreatedTag(tt["existing_user_and_repo"].user, tt["existing_user_and_repo"].requestBody, tt["existing_user_and_repo"].repo)
}

func TestGetARepoRecommendation(t *testing.T) {
	tt := map[string]struct {
		user           string
		repo           int64
		responseStatus int
		responseBody   string
	}{
		"correct_user":         {"joaopmgd", 10866521, http.StatusOK, ""},
		"no_user_found":        {"joaopmgdjoaopmgdjoaopmgdjoaopmgd", 10866521, http.StatusNotFound, "{\"error\":\"User not found\"}"},
		"no_user":              {"", 0, http.StatusNotFound, "{\"error\":\"User not found\"}"},
		"no_repo_found":        {"joaopmgd", 0, http.StatusNotFound, "{\"error\":\"Repository not found 0\"}"},
		"recommendation_found": {"joaopmgd", 10866521, http.StatusOK, ""},
	}
	for testName, tc := range tt {

		req, err := http.NewRequest("GET", "/repos/{user}/starred/{repo}/recommendation", nil)
		if err != nil {
			t.Fatal(err)
		}
		req = mux.SetURLVars(req, map[string]string{
			"user": tc.user,
			"repo": strconv.FormatInt(tc.repo, 10),
		})

		response := executeRequestTest(req, a.GetARepoRecommendation)

		if response.Code != tc.responseStatus ||
			(response.Code != http.StatusOK && response.Body.String() != tc.responseBody) {
			t.Errorf("\nHandler %s\nFor user %s\nGot Status %v and Body %s\nWant Status %v and Body %s",
				testName, tc.user, response.Code, response.Body.String(), tc.responseStatus, tc.responseBody)
		}
	}
}

func executeRequestTest(req *http.Request, appHandler http.HandlerFunc) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(appHandler)
	handler.ServeHTTP(rr, req)
	return rr
}

// Remove tag created to end the test
func removeCreatedTag(user, body string, repo int64) {
	req, _ := http.NewRequest("POST", "/repos/{user}/starred/{repo}", strings.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{
		"user": user,
		"repo": strconv.FormatInt(repo, 10),
	})
	executeRequestTest(req, a.DeleteTagStarredRepo)
}
