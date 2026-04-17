.PHONY: test-integrations

test-integrations:
	docker-compose -f docker-compose.test.yaml up -d
	sleep 5
	POSTGRES_PORT_EXTERNAL=5436 go run cmd/migrate/main.go -command up
	go test -tags=integration ./...
	docker-compose -f docker-compose.test.yaml down

.PHONY: migration-up

migration-up:
	go run cmd/migrate/main.go -command up
