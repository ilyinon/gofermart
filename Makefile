APP=gophermart

DB_URI=postgres://user:pass@localhost:5432/gofermart?sslmode=disable

.PHONY: build run test up down migrate migrate-down

build:
	go build -o bin/$(APP) ./cmd/gophermart/main.go

run:
	go run ./cmd/gophermart/main.go

test:
	go test ./... -cover

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

db:
	docker exec -it gophermart-db psql -U user -d gofermart

migrate:
	goose -dir migrations postgres "$(DB_URI)" up

migrate-down:
	goose -dir migrations postgres "$(DB_URI)" down
