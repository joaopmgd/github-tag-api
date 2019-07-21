package model

// RequestError will detail the error encountered with the github API
type RequestError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

// StarredRepo is the starred repo complete data
type StarredRepo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Language    string `json:"language"`
	RequestError
}
