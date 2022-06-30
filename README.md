# RUNNING THE API

### Required Installations:

- Go (used 1.18)
- PostgreSQL (used v14)
- [golang-migrate](https://github.com/golang-migrate/migrate) (used v4.14.1)

_More information [here](https://pkg.go.dev/github.com/golang-migrate/migrate/cli#section-readme) on installing golang-migrate on different platforms_

### Instructions:

1. Create a `cards` database.
2. Install `citext` and `uuid-ossp` extensions.
3. Create a `cards` role with password `passwOrd`.

_You can skip steps one to three if you already have an existing role and db (that has the extensions isntalled) you want to use._

4. Populate `CARDS_DB_DSN` variable in `.envrc`. It has a default value of `postgres://cards:password@localhost/cards` with `cards` as db and role name and `passwOrd` as password. If using your own db and role, feel free to change the default value. (PostgreSQL DSN format: postgres://<username>:<password>@localhost/<dbname>)
5. Run `make db/migrations/up` from the root of the app, which runs the migrations inside the `migrations` folder.
6. Run `go mod download` to download Go module dependencies.
7. Run `make run/api` from the root of the app to start the API.

_If on Linux, specifically Ubuntu, you can run the script `setup.sh` which does everything listed above. However please make sure that you have Go, PostgreSQL and golang-migrate installed before running the script._

# RUNNING THE TESTS

- The handler tests can be run by moving into the `cmd/api` folder and then running `go test`.
- The integration test, however, has a little bit of setup needed. If you look at the function `newTestDB` in `testutils_test.go`, it expects a `test` database and a `test` role with password `passwOrd`. If you have an existing db and role you want to use for testing, edit the second argument of `sql.Open` in the mentioned `newTestDB` function. If not, then setup the new database using the instructions above. Don't forget the extensions.

# ROUTES

| Method | Path          | Description               | Payload                   | Response                                              |
| ------ | ------------- | ------------------------- | ------------------------- | ----------------------------------------------------- |
| GET    | /v1/decks/:id | Get information on a deck | NONE                      | JSON (deck_id, remaining, shuffled, array of cards\*) |
| POST   | /v1/decks     | Create a deck             | JSON (shuffled and cards) | JSON (deck_id, remaining, shuffled)                   |
| PUT    | /v1/decks/:id | Draw cards from deck      | JSON (count)              | JSON (array of cards\*)                               |

\*Each card is a JSON object with value, suit, and code fields

# ADDITIONAL NOTES

(in addition to specs listed [here](https://toggl.notion.site/Toggl-Backend-Unattended-Programming-Test-015a95428b044b4398ba62ccc72a007e))

### POST /v1/decks

- Default value of `shuffled` is `false`. If JSON payload doesn't have the field `shuffled`, then it defaults to `false`.
- Default value of cards is a full deck, which means that a missing `cards` field or an empty array value for `cards`, will create a deck with 52 cards.

### GET /v1/decks/:id

- If a deck returned has been fully dealt, `remaining` will be `0` and there will be no `cards` field returned.
