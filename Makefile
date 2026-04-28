.PHONY: build up down ps fmt lint test frontend-fmt frontend-lint frontend-test backend-fmt backend-lint backend-test

# Docker operations (Containers)
build:
	docker compose build

up:
	docker compose up --build

down:
	docker compose down

ps:
	docker compose ps

# Local development tasks (Run on host machine)
fmt: frontend-fmt backend-fmt
lint: frontend-lint backend-lint
test: frontend-test backend-test

# Frontend local tasks
frontend-fmt:
	cd frontend && npm run fmt

frontend-lint:
	cd frontend && npm run lint

frontend-test:
	cd frontend && npm run test

# Backend local tasks
backend-fmt:
	cd backend && go fmt ./...

backend-lint:
	cd backend && go vet ./...

backend-test:
	cd backend && go test ./...

