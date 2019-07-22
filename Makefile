run:
	HOST=:8080 \
	GITHUB_PROPERTIES_ENDPOINT='https://api.github.com' \
	GITHUB_USER_STARRED='/users/{{ .user }}/starred' \
	GITHUB_HEALTH_STATUS='https://www.githubstatus.com/api/v2/status.json' \
	DB_HOST=localhost \
	DB_PORT=5432 \
	DB_USER=root \
	DB_DB_NAME=github \
	DB_PASSWORD=123456 \
	go run main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o github-tag-api-linux main.go

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o github-tag-api-mac main.go

docker-build:
	docker build -t github-tag-api .

docker-run:
	docker run -p 8080:8080 github-tag-api