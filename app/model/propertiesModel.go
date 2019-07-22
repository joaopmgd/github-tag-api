package model

// RequestError will detail the error encountered with the github API
type RequestError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

// StarredRepoRequest is the starred repo complete data
type StarredRepoRequest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Language    string `json:"language"`
	RequestError
}

// StarredRepoTags is the starred repo complete data with tags that were added
type StarredRepoTags struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Language    string   `json:"language"`
	Tags        []string `json:"tags"`
}

// TagRequestUpdate is the body from the tag request POST
type TagRequestUpdate struct {
	TagName string `json:"tag"`
}

// ResponseOK responds a message when the request is ok
type ResponseOK struct {
	Message string `json:"Message"`
}

// RecommendedTags list of Recommended tags for a certain language
type RecommendedTags struct {
	Recommended []string `json:"recommended"`
}

// StarredRepoTagsResponse has the pagination and RESTful data added to the response list
type StarredRepoTagsResponse struct {
	StarredRepos         []StarredRepoTags `json:"starred_repos"`
	PageNumber           int               `json:"page_number"`
	PageSize             int               `json:"page_size"`
	PropertiesTotalCount int               `json:"properties_total_count"`
}

// GithubHealthStatus stores the health status from github
type GithubHealthStatus struct {
	Status GithubStatus `json:"status"`
}

// GithubStatus status data from github
type GithubStatus struct {
	Indicator   string `json:"indicator"`
	Description string `json:"description"`
}

// AppHealthStatus shows the app health status
type AppHealthStatus struct {
	Status       string         `json:"status"`
	GithubStatus GithubStatus   `json:"github"`
	Database     DatabaseStatus `json:"database"`
}

// DatabaseStatus detaisl the database health status
type DatabaseStatus struct {
	Status string `json:"status"`
}
