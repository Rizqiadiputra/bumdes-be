run:
	go run ./cmd/api

migrate:
	go run ./cmd/migrate

seed:
	go run ./cmd/seed

billing-cron:
	go run ./cmd/billing-cron

build:
	go build -o bin/api ./cmd/api
	go build -o bin/migrate ./cmd/migrate
	go build -o bin/seed ./cmd/seed
	go build -o bin/billing-cron ./cmd/billing-cron

tidy:
	go mod tidy

swag:
	swag init -g cmd/api/main.go -o docs
