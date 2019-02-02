all: build

.PHONY build:
build:
	go build -o bin/server

.PHONY server:
server: build
	./bin/server
