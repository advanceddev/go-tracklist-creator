APP_CMD_PATH=./cmd/go-tracklist-creator
BINARY_NAME=go-tracklist-creator

# Цели
all: build

build:
		go build -o $(BINARY_NAME) $(APP_CMD_PATH)

dev:
		go run $(APP_CMD_PATH)

test:
		go test ./...

clean:
		rm -f $(BINARY_NAME)

.PHONY: all run test clean build