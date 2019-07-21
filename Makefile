run:
	HOST=:3000 \
	GITHUB_USER_STARRED='/users/{{ .user }}/starred' \
	GITHUB_PROPERTIES_ENDPOINT='https://api.github.com' \
	go run main.go
