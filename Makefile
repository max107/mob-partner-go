UNAME_S := $(shell uname -s)

all: build

.PHONY build:
build:
	GOOS=linux GOARCH=amd64 go build -o bin/linux_server
	GOOS=darwin GOARCH=amd64 go build -o bin/mac_server

.PHONY server:
server: build
    ifeq ($(UNAME_S),Linux)
	./bin/linux_server ./config.json
    endif
    ifeq ($(UNAME_S),Darwin)
	./bin/mac_server ./config.json
    endif

.PHONY deploy:
deploy: build
	rsync -arcvp ./bin/linux_server user@test-target:/home/user/go/
	ssh user@test-target "sudo systemctl restart goserver"
