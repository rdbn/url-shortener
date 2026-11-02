.PHONY: lint test-integrations migrate-docker-up migrate-docker-down

lint:
	golangci-lint run ./...

test-integrations:
	docker compose -f docker-compose.test.yaml up -d postgres_test redis_test
	sleep 5
	docker compose -f docker-compose.test.yaml run --rm migrate_test
	POSTGRES_HOST=localhost go test -tags=integration ./...
	docker compose -f docker-compose.test.yaml down

migrate-docker-up:
	docker compose exec app sh -ec 'migrate -path ./migrations -database "postgres://$${POSTGRES_USER}:$${POSTGRES_PASSWORD}@postgres:5432/$${POSTGRES_DB}?sslmode=$${POSTGRES_SSLMODE}" up'

migrate-docker-down:
	docker compose exec app sh -ec 'migrate -path ./migrations -database "postgres://$${POSTGRES_USER}:$${POSTGRES_PASSWORD}@postgres:5432/$${POSTGRES_DB}?sslmode=$${POSTGRES_SSLMODE}" down 1'
