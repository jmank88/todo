.PHONY: all
all: clean test build

.PHONY: clean
clean:
	go clean

# Runs all unit tests.
.PHONY: test
test:
	go test $(shell go list ./... | grep -v /vendor/)

# Builds the server and cli binaries.
.PHONY: build
build:
	go build .
	go build ./client/cli

# Runs datastore integration tests
.PHONY: integration-tests
integration-tests:
	docker-compose build
	docker-compose start postgres

	docker-compose run todo go test -tags=integration ./datastore -host=postgres://postgres:postgres@postgres:5432?sslmode=disable

	docker-compose stop