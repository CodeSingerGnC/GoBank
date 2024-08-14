#!/bin/sh

set -e 

echo "run db migration"
/app/migrate -path /app/migration -database "$DB_CONNECTION_SOURCE" -verbose up
echo "db migration success"

echo "start app"
exec "$@"