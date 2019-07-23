# Golang Github Tag API

This RESTful API will recover a user Starred repositories, add some tags to it so it can classify better the selection. With all this data it will be possible to recover the repos with the same tags e and even recomend some tags to be added.

## Prerequisites

This project was built in Go version go1.12, so it must be a version similar with that.

NOTE: dependency management was done using Dep (dependency management tool for Go)

## Running

To run this project just enter the command bellow:

Running in the terminal if you have golang installed:
```
make run
```

Building and Creating a Container:
```
make docker-build
```

Running the Container locally:
```
make docker-run
```

## Requests 

### GET repos/{username}/starred?tag={tag}

- To recover all starred repos by an user, the GET request will only need an URL parameter for the username. If a tag is passed in the query params the search will return starred repos that were tagged with that search information

### GET /repos/{user}/starred/{repo}/recommendation

- To recover all starred repos by an user, the GET request will need an URL parameter called for the username and for the repo that should be recommendated

### POST /repos/{user}/starred/{repo}

- To add a new tag for a repo
- The body for the post request should be a JSON as:
{
	"tag": "test"
}

### DELETE /repos/{user}/starred/{repo}

- To add a delete tag for a repo
- The body for the post request should be a JSON as:
{
	"tag": "test"
}



Pagination is implemented, the default Response has offset=0 and limit=10 for starred repos, if there is a need to change that just run the request with the below for example:
```
/repos/{username}/starred?offset=0&limit=10
```


## Running the tests

The file main_test.go has the integration tests, which is required a connection to the db and Github

To run the tests locally just execute:

```
go test
```


## Running on Docker

The endpoint and the Port are external varibles setup in the Dockerfile, if change is needed it is ok, but they are required.

Thanks ;D