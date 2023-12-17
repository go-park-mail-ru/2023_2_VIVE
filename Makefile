.PHONY: migrate
migrate:
	dotenv -- tern migrate -m ./deploy/migrations/hnh

.PHONY: rollback
rollback:
	dotenv -- tern migrate -m ./deploy/migrations/hnh -d -1

.PHONY: run
run:
	go run ./cmd/HnH/HnH.go

.PHONY: search_postgres
search_postgres: deploy/Dockerfile
	docker build -t search_postgres ./deploy/

/PHONY: compose
compose: search_postgres
	dotenv -- docker compose -f ./deploy/docker-compose.yaml up -d

.PHONY: test
test:
	go test -coverpkg=./... -coverprofile=c.out ./...

.PHONY: cover
cover: test
	go tool cover -func=c.out
