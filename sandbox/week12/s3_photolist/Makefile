COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

export GO111MODULE=on

.PHONY: build
build: assets
	@echo "-- building binary"
	go build \
		-ldflags "-X main.buildHash=${COMMIT} -X main.buildTime=${BUILD_TIME}" \
		-o ./bin/photolist \
		./cmd/photolist

.PHONY: assets
assets: 
	@echo "-- generatiog assets"
	go run pkg/assets/assets_gen.go

.PHONY: gqlgen
gqlgen: gen

.PHONY: gen
gen: 
	@echo "-- generatiog graphql files"
	go run github.com/99designs/gqlgen generate --verbose --config ./configs/gqlgen.yml

.PHONY: dev
dev: 
	@echo "-- starting air wrapper"
	air -c ./configs/air.conf

.PHONY: docker
docker: 
	@echo "-- building docker container - multistage"
	docker build -f build/Dockerfile.Multistage -t photolist .

.PHONY: docker_run
docker_run: 
	@echo "-- starting docker container"
	docker run -it -p 8080:8080 photolist

.PHONY: dc
dc: docker_compose

.PHONY: docker_compose
docker_compose: 
	@echo "-- starting docker compose"
	docker-compose -f ./deployments/docker-compose.yml up

.PHONY: dcb
dcb: 
	@echo "-- starting docker compose with build"
	docker-compose -f ./deployments/docker-compose.yml up --build
