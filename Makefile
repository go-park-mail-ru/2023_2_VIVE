-include .env
include versions.env

.PHONY: create-migration
create-migration:
	tern new -m deploy/migrations/hnh/ $(name)

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

.PHONY: compose-up
compose-up:
	dotenv -- docker compose -f ./deploy/docker-compose.yaml up -d

.PHONY: compose-down
compose-down:
	dotenv -- docker compose -f ./deploy/docker-compose.yaml down

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -coverpkg=./... -coverprofile=c.out ./...

.PHONY: cover
cover: test
	go tool cover -func=c.out

search_db_image = vovchenskiy/hnh:search_db-v${SEARCH_DB_VERSION}

.PHONY: build-hnh_search-db
build-hnh-search-db: deploy/Dockerfile
	docker build -t ${search_db_image} ./deploy/

main_service_image = vovchenskiy/hnh:main-v${MAIN_SRVC_VERSION}

.PHONY: build-hnh-main
build-hnh-main: Dockerfile
	docker build . --file Dockerfile --tag ${main_service_image}

auth_service_image = vovchenskiy/hnh:auth-v${MAIN_SRVC_VERSION}

.PHONY: build-hnh-auth
build-hnh-auth: services/auth/Dockerfile
	docker build --file ./services/auth/Dockerfile --tag ${auth_service_image} .

search_service_image = vovchenskiy/hnh:search-v${MAIN_SRVC_VERSION}

.PHONY: build-hnh-search
build-hnh-search: services/searchEngineService/Dockerfile
	docker build --file ./services/searchEngineService/Dockerfile --tag ${search_service_image} .

notifications_service_image = vovchenskiy/hnh:notifications-v${MAIN_SRVC_VERSION}

.PHONY: build-hnh-notifications
build-hnh-notifications: services/notifications/Dockerfile
	docker build --file ./services/notifications/Dockerfile --tag ${notifications_service_image} .

csat_service_image = vovchenskiy/hnh:csat-v${MAIN_SRVC_VERSION}

.PHONY: build-hnh-csat
build-hnh-csat: services/csat/Dockerfile
	docker build --file ./services/csat/Dockerfile --tag ${csat_service_image} .

.PHONY: build-hnh
build-hnh: build-hnh-main build-hnh-auth build-hnh-search build-hnh-notifications build-hnh-csat build-hnh_search-db

.PHONY: push-images
push-images:
	docker push ${main_service_image}
	docker push ${auth_service_image}
	docker push ${search_service_image}
	docker push ${notifications_service_image}
	docker push ${csat_service_image}
	docker push ${search_db_image}

.PHONY: pull-images
pull-images:
	dotenv -e .env -e versions.env -- docker compose -f deploy/docker-compose.yaml pull
