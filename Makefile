run/api:
	@go run ./cmd/api

psql:
	psql ${CARDS_DB_DSN}

db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${CARDS_DB_DSN} up

db/migrations/down:
	@echo 'Running down migrations...'
	migrate -path ./migrations -database ${CARDS_DB_DSN} down
