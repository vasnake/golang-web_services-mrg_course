COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

all: run

.PHONY: gqlgen_init
gqlgen_init: 
	@echo "-- init gqlgen"
	go run github.com/99designs/gqlgen init

.PHONY: gen
gen: 
	@echo "-- running gqlgen"
	go run github.com/99designs/gqlgen generate --verbose

.PHONY: run
run:
	@echo "-- starting app"
	go run ./server/server.go

