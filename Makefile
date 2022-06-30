include .envrc

## help: print help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm: 
	@echo -n 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

.PHONY: run/api
## run/api: run the cmd/api application
run/api:
	@go run ./cmd/api -db-dsn=${CARDS_DB_DSN}

.PHONY: build/api
## build/api: build the cmd/api application
build/api:
	@echo 'BUilding cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api

# ==================================================================================== #
# DB
# ==================================================================================== #

.PHONY: db/psql
## db/psql: connect to the database using psql
db/psql:
	psql ${CARDS_DB_DSN}

.PHONY: db/migrations/up
## db/migrations/up: apply all up database migrations
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations/api -database ${CARDS_DB_DSN} up

.PHONY: db/migrations/down
## db/migrations/down: apply all down database migrations
db/migrations/down: confirm
	@echo 'Running down migrations...'
	migrate -path ./migrations/api -database ${CARDS_DB_DSN} down

# ==================================================================================== #
# AUDIT
# ==================================================================================== #
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

