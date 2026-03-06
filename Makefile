APP_NAME=devpulse-api

run:
	go run ./cmd/api

compose-up:
	docker compose up -d

compose-down:
	docker compose down

fmt:
	go fmt ./...

tidy:
	go mod tidy