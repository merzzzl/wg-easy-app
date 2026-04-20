.PHONY: install build build-backend build-frontend run-backend run-frontend compose-up compose-down compose-logs bootstrap-deploy

install:
	npm --prefix frontend install
	cd backend && go mod tidy

build: build-backend build-frontend

build-backend:
	cd backend && go build ./...

build-frontend:
	npm --prefix frontend run build

run-backend:
	cd backend && go run ./cmd/app

run-frontend:
	npm --prefix frontend run dev

compose-up:
	docker compose -f compose/compose.yaml up -d

compose-down:
	docker compose -f compose/compose.yaml down

compose-logs:
	docker compose -f compose/compose.yaml logs -f

bootstrap-deploy:
	./scripts/bootstrap-deploy.sh
