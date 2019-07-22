# Based on the latest golang image from Docker
FROM golang:latest AS builder

ADD . /go/src/github.com/joaopmgd/github-tag-api
WORKDIR /go/src/github.com/joaopmgd/github-tag-api

# Building the Go executable for linux
RUN GOOS=linux GOARCH=386 go build -o /app -i main.go

########################################################

FROM scratch

ENV HOST=:8080 \
    GITHUB_PROPERTIES_ENDPOINT='https://api.github.com' \
    GITHUB_USER_STARRED='/users/{{ .user }}/starred' \
    GITHUB_HEALTH_STATUS='https://www.githubstatus.com/api/v2/status.json'\
    DB_HOST=localhost \
    DB_PORT=5432 \
    DB_USER=root \
    DB_DB_NAME=github \
    DB_PASSWORD=123456

COPY --from=builder /app ./
WORKDIR /

EXPOSE 8080

ENTRYPOINT ["./app"]