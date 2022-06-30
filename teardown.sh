make db/migrations/down
sudo -i -u postgres psql -c "DROP DATABASE cards"
sudo -i -u postgres psql -d cards -c 'DROP ROLE cards'

