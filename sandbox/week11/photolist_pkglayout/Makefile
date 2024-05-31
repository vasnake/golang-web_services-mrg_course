COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

all: dev

.PHONY: dev
dev: 
	@echo "-- starting app in dev mode with air"
	air -c ./configs/air.conf

.PHONY: emb
emb: 
	@echo "-- running templates embedding"
	go generate --tags=dev

.PHONY: build
build: emb
	@echo "-- compiling binary"
	go build \
		-ldflags "-X main.buildHash=${COMMIT} -X main.buildTime=${BUILD_TIME}" \
		-o ./bin/photolist.exe \
		./cmd/photolist

.PHONY: assets
assets: 
	@echo "-- generatiog assets"
	go run pkg/assets/assets_gen.go

.PHONY: gen
gen: 
	@echo "-- generatiog graphql files"
	go run github.com/99designs/gqlgen generate --verbose --config ./configs/gqlgen.yml

.PHONY: dc
dc: 
	@echo "-- starting docker compose"
	docker-compose -f ./dev/docker-compose.yaml up
