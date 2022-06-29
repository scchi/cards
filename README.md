Requirements to run the program:

- Go
- PostgreSQL
- golang-migrate (v4.14.1) https://github.com/golang-migrate/migrate

Instructions for running the program:

1. Create a database
2. Install citext and uuid-ossp extensions
3. Create role for the app

You can skip steps one to three if you already have an existing role and db you want to use.

4. By default, app assumes there is an environment variable named CARDS_DB_DSN to connect to the database. If you don't want to define a CARDS_DB_DSN environment variable, you can provide it as db-dsn argument when running go run command. (PostgreSQL DSN format: postgres://<username>:<password>@localhost/<dbname>)
5. Run migrations (migrate -path=../migrations -database=$CARDS_DB_DSN up)
6. Run `go run ./cmd/api` from the root of the app

If running Ubuntu, you can run the script setup.sh. However please make sure that you have Go, PostgreSQL and golang-migrate installed before running the script.
