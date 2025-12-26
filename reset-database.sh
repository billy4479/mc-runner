#!/bin/sh

set -euo pipefail

DB_PATH=$1

rm -v "$DB_PATH"
migrate -path ./db/migrations -database "sqlite3://$DB_PATH" up

