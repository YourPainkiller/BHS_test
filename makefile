APP=app.exe
APP_EXECUTABLE="./cmd/app/$(APP)"
SHELL := /bin/bash 
PG_URL="postgres://postgres:qwe@localhost:5434/store?sslmode=disable"


up:
	make tidy
	make compose-up
	make goose-up
	make run

check-quality: 
	make fmt
	make vet

cover:
	go test -cover ./internal/usecase
	
vet: 
	go vet ./...

fmt: 
	go fmt ./...

tidy: 
#	go get -u
	go mod tidy

build:
	go build -o $(APP_EXECUTABLE) ./cmd/app/main.go
	@echo "Build passed"

run: 
	make build
	$(APP_EXECUTABLE)

gen-swag:
	swag init -g cmd/app/main.go

# ---------------------------
# Запуск базы данных в Docker
# ---------------------------

compose-up:
	docker-compose up -d 

compose-down:
	docker-compose down

compose-stop:
	docker-compose stop 

compose-start:
	docker-compose start 

compose-ps:
	docker-compose ps 

# ---------------------------
# Запуск миграций через Goose
# ---------------------------

goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

goose-add:
	goose -dir ./migrations postgres $(PG_URL) create rename_me sql

goose-up:
	goose -dir ./migrations postgres $(PG_URL) up

goose-down:
	goose -dir ./migrations postgres $(PG_URL) down

goose-status:
	goose -dir ./migrations postgres $(PG_URL) status
