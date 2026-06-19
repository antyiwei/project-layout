.PHONY: build run test lint clean

BINARY := bin/server

build:
	go build -o $(BINARY) ./cmd/server

run: build
	./$(BINARY)

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -rf bin/
