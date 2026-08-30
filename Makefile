include .env
export

export PROJECT_ROOT=$(shell pwd)

COMPOSE_DEV = docker compose \
	--env-file .env \
	-f deployments/docker/compose.yaml \
	-f deployments/docker/compose.dev.yaml




env-up:
	@$(COMPOSE_DEV) up -d postgres redis prometheus grafana loki alloy

env-down:
	@$(COMPOSE_DEV) down

env-cleanup:
	@read -p "Clear all volume files in the environment? Risk of data loss. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		$(COMPOSE_DEV) down -v; \
		echo "Environment files have been cleared. "; \
	else \
		echo "Environment cleanup cancelled. "; \
	fi

ps:
	@docker ps -a




migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Missing seq parameter selection. Example: make migrate-create seq=init" \
		exit 1; \
	fi; \
	$(COMPOSE_DEV) run --rm go-migrations \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@$(MAKE) migrate-action action=up

migrate-down:
	@$(MAKE) migrate-action action=down 

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Missing action parameter selection. Example: make migrate-action  action=up" \
		exit 1; \
	fi; \
	$(COMPOSE_DEV) run  --rm go-migrations \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"



auth-run: 
	@go run ./cmd/auth-service/main.go

auth-test:
	@go test ./internal/service/auth/... -v -coverprofile=coverage.out

bash-test:
	@./bash_tests/test.sh



auth-build:
	@$(COMPOSE_DEV) build auth-service

auth-up-b:
	@$(COMPOSE_DEV) up -d --build auth-service

auth-down:
	@$(COMPOSE_DEV) stop auth-service

auth-logs:
	@$(COMPOSE_DEV) logs -f auth-service

auth-up:
	@$(COMPOSE_DEV) up -d auth-service


swagger-gen:
	@$(COMPOSE_DEV) run --rm swagger \
		init \
		-g cmd/auth-service/main.go \
		-o docs \
		--parseInternal \
		--parseDependency 