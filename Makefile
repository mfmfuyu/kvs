SERVER_SRC=server/main.go
CLI_SRC=cli/main.go

.PHONY: server build-server cli build-cli

server:
	go run $(SERVER_SRC)

build-server:
	go build -o dist/server $(SERVER_SRC)

cli:
	go run $(CLI_SRC)

build-cli:
	go build -o dist/cli $(CLI_SRC)

build: build-server build-cli
