# TODO: check ${OAUTH_APP_ID} and ${OAUTH_APP_SECRET}

# https://www.gnu.org/software/make/manual/html_node/Values.html
.EXPORT_ALL_VARIABLES:
# BUILD_TIME != date -u '+%Y-%m-%d_%H:%M:%S'UTC
# BUILD_HASH != cd ../..; git rev-parse --short HEAD
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d')
BUILD_HASH ?= $(shell cd ../..; git rev-parse --short HEAD)
GOWORK=off
# export GO111MODULE=on

# build -------------------------------------------------------

.PHONY: build
build:
build: assets
	@echo "-- building binaries"
	@echo check_this_line_for_build_vars go build -ldflags -X main.buildHash=${BUILD_HASH} -X main.buildTime=${BUILD_TIME}
	go build --tags=prod \
		-ldflags "-X main.buildHash=${BUILD_HASH} -X main.buildTime=${BUILD_TIME}" \
		-o /tmp/photolist/bin/photolist \
		./cmd/photolist
	go build \
		-ldflags "-X main.buildHash=${BUILD_HASH} -X main.buildTime=${BUILD_TIME}" \
		-o /tmp/photolist/bin/photoauth \
		./cmd/photoauth
	go build \
		-ldflags "-X main.buildHash=${BUILD_HASH} -X main.buildTime=${BUILD_TIME}" \
		-o /tmp/photolist/bin/auth \
		./cmd/auth

.PHONY: assets
assets: clean
	@echo "-- generatiog assets"
	go run cmd/assets_gen/assets_gen.go

.PHONY: clean
clean: 
	@echo "-- drop generated files"
	rm -v -- pkg/assets/assets_vfsdata.go || echo file not exists

.PHONY: graphql
graphql: 
	@echo "-- generatiog graphql files"
	go run github.com/99designs/gqlgen -c ./configs/gqlgen.yml

## wtf, where is grpc proto gen?

# CI/CD -------------------------------------------------------

# you should create file `./secrets.env` with content like this:
# OAUTH_APP_ID=Ov***JF
# OAUTH_APP_SECRET=ad***60
.PHONY: up
up: compose_build
	@echo "-- starting docker compose"
	# check: docker compose -f ./deployments/docker-compose.yml --env-file ./secrets.env config
	docker compose -f ./deployments/docker-compose.yml --env-file ./secrets.env up&

.PHONY: down
down: 
	@echo "-- starting docker compose"
	docker-compose -f ./deployments/docker-compose.yml down

# Build images before starting containers
.PHONY: compose_build
compose_build:
	@echo "-- starting docker compose with build"
	docker-compose -f ./deployments/docker-compose.yml \
		build --build-arg BUILD_HASH=${BUILD_HASH} --build-arg BUILD_TIME=${BUILD_TIME}

# build app binaries container
.PHONY: docker_build
docker_build:
	@echo "-- building docker container - multistage, CURDIR: $(CURDIR)"
	pwd
	# env vars delivery option, see dockerfile
	# echo "BUILD_HASH=${BUILD_HASH}" > ./build.env
	# echo "BUILD_TIME=${BUILD_TIME}" >> ./build.env
	docker build \
		--build-arg BUILD_HASH=${BUILD_HASH} --build-arg BUILD_TIME=${BUILD_TIME} \
		-f ./build/Dockerfile.Multistage -t photolist:latest .
