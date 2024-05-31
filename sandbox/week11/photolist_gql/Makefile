COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

all: run

.PHONY: air
air: install_air dev

.PHONY: install_air
install_air: 
	@echo "-- installing air"
	go install github.com/cosmtrek/air@latest
	air -v

.PHONY: install_vfsgen
install_vfsgen: 
	@echo "-- installing vfsgen"
	go get github.com/shurcooL/vfsgen

.PHONY: dev
dev: 
	@echo "-- starting app in dev mode with air"
	air -c air.conf

.PHONY: emb
emb: 
	@echo "-- running templates embedding"
	go generate --tags=dev

.PHONY: build
build: emb
	@echo "-- compiling binary"
	go build -o ./photolist.exe .

.PHONY: run
run: build
	@echo "-- starting app"
	./photolist.exe

.PHONY: dc
dc: 
	@echo "-- starting docker compose"
	docker-compose up

.PHONY: gen
gen: 
	@echo "-- running gqlgen"
	go run github.com/99designs/gqlgen generate --verbose

