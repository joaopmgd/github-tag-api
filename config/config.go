package config

import (
	"bytes"
	"html/template"
	"os"
)

// Config will setup the Endpoints, the sources that will be requested, Log and Memory
type Config struct {
	Endpoints *Endpoint
	Log       *StandardLogger
}

// Endpoint for the future Requests
type Endpoint struct {
	GithubURL         string
	GithubUserStarred string
}

// GetConfig will setup the config struct for the app to run
func GetConfig() *Config {
	return &Config{
		Endpoints: &Endpoint{
			GithubURL:         os.Getenv("GITHUB_PROPERTIES_ENDPOINT"),
			GithubUserStarred: os.Getenv("GITHUB_USER_STARRED"),
		},
		Log: NewLogger(),
	}
}

// GetStarredReposURL creates the starred repos url
func (c *Config) GetStarredReposURL(vars map[string]string) (string, error) {
	templateURL := c.Endpoints.GithubURL + c.Endpoints.GithubUserStarred
	t, err := template.New("URL").Parse(templateURL)
	if err != nil {
		c.Log.CreatingRestTemplateError(err.Error())
		return "", err
	}
	var URL bytes.Buffer
	err = t.Execute(&URL, vars)
	if err != nil {
		c.Log.ExecutinRestTemplateError(err.Error())
		return "", err
	}
	return URL.String(), nil
}
