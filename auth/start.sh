#!/bin/sh

set -e

echo "run db migrations . . . "
source .env
/app/migrate -path /app/migrations -database "$PG_TEST_URL" -verbose up

echo "start the app"
exec "$@"
