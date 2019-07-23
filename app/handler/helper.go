package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/joaopmgd/github-tag-api/app/model"
	"github.com/joaopmgd/github-tag-api/config"
)

func addLanguage(language string, tags []string) []string {
	if tags == nil {
		tags = []string{}
	}
	for _, tag := range tags {
		if tag == language {
			return tags
		}
	}
	if len(language) > 0 {
		return append(tags, language)
	}
	return tags
}

func repoHasTag(tags []string, selectedTag string) bool {
	if selectedTag != "" {
		for _, tag := range tags {
			if strings.Contains(tag, selectedTag) {
				return true
			}
		}
	}
	return false
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
