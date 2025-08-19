#!/bin/sh

# Parse .env - https://gist.github.com/judy2k/7656bfe3b322d669ef75364a46327836
export $(egrep -v '^#' .env | xargs)

rm -v "$DB_PATH"
migrate -path ./db/migrations -database "sqlite3://$DB_PATH" up

