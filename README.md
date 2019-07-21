# Golang Github Tag API

This RESTful API will recover a user Starred repositories, add some tags to it so it can classify better the selection. With all this data it will be possible to recover the repos with the same tags e and even recomend some tags to be added.

### Prerequisites

This project was built in Go version go1.12, so it must be a version similar with that.

NOTE: dependency management was done using Dep (dependency management tool for Go)

### Running

To run this project just enter the command bellow:

Running in the terminal:
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


## Running the tests

To run the tests just execute:

```
go test
```


## Running on Docker

The endpoint and the Port are external varibles setup in the Dockerfile, if change is needed it is ok, but they are required.

Thanks ;D