all: build

.PHONY build:
build:
	GOOS=linux GOARCH=amd64 go build -o bin/linux_server
	GOOS=darwin GOARCH=amd64 go build -o bin/mac_server

.PHONY server:
server: build
	UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        ./bin/linux_server
    endif
    ifeq ($(UNAME_S),Darwin)
        ./bin/mac_server
    endif
