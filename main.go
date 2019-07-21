package main

import (
	"os"

	"github.com/github-tag-api/app"
	"github.com/github-tag-api/config"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config := config.GetConfig()
	app := &app.App{}
	app.Initialize(config)
	app.Run(os.Getenv("HOST"))
}
