GOCMD = go
GOBUILD = $(GOCMD) build
BINARY_NAME = go-tcp-client-agent

run:
	GO111MODULE=on $(GOCMD) run .

build:
	mkdir -p ./build
	GO111MODULE=on $(GOBUILD) -o ./build/$(BINARY_NAME) .

clean:
	rm -rf ./build/*

.PHONY: run build clean
