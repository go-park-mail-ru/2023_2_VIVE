include .env

.PHONY: migrate
migrate:
	dotenv -- tern migrate -m ./deploy/migrations/hnh

.PHONY: rollback
rollback:
	dotenv -- tern migrate -m ./deploy/migrations/hnh -d -1

.PHONY: build
build:
	go build -o ./bin/hnh ./cmd/HnH/HnH.go

.PHONY: run
run:
	go run ./cmd/HnH/HnH.go

.PHONY: search_postgres
search_postgres: /deploy/Dockerfile
	docker build -t search_postgres ./deploy/

.PHONY: compose
compose:
	dotenv -- docker compose -f ./deploy/docker-compose.yaml up -d

.PHONY: test
test:
	go test -coverpkg=./... -coverprofile=c.out ./...

.PHONY: cover
cover: test
	go tool cover -func=c.out

.PHONY: build-hnh-main
build-hnh-main: Dockerfile
	docker build . --file Dockerfile --tag hnh_main:${HNH_VERSION}

.PHONY: build-hnh-auth
build-hnh-auth:
	docker build --file ./services/auth/Dockerfile --tag hnh_auth:${AUTH_SRVC_VERSION} .

.PHONY: build-hnh-search
build-hnh-search:
	docker build --file ./services/searchEngineService/Dockerfile --tag hnh_search:${SEARCH_SRVC_VERSION} .

.PHONY: build-hnh-notifications
build-hnh-notifications:
	docker build --file ./services/notifications/Dockerfile --tag hnh_notifications:${NOTIFICATION_SRVC_VERSION} .

.PHONY: build-hnh-csat
build-hnh-csat:
	docker build --file ./services/csat/Dockerfile --tag hnh_csat:${CSAT_SRVC_VERSION} .

.PHONY: build-hnh
build-hnh: build-hnh-main build-hnh-auth build-hnh-search build-hnh-notifications build-hnh-csat

