read -p "Enter password for DB user: " DB_PASSWORD
echo "CARDS_DB_DSN='postgres://cards:${DB_PASSWORD}@localhost/cards'" > .envrc

curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate.linux-amd64 /usr/local/bin/migrate

sudo apt install postgresql
sudo -i -u postgres psql -c "CREATE DATABASE cards"
sudo -i -u postgres psql -d cards -c "CREATE EXTENSION IF NOT EXISTS citext"
sudo -i -u postgres psql -d cards -c 'CREATE EXTENSION IF NOT EXISTS "uuid-ossp"'
sudo -i -u postgres psql -d cards -c "CREATE ROLE cards WITH LOGIN PASSWORD '${DB_PASSWORD}'"

make db/migrations/up

go mod download
make run/api