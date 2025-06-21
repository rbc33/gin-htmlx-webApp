#! /usr/bin/env bash
apk add git netcat-openbsd
set -euo pipefail
git config --global --add safe.directory '*'

# Wait for MySQL to be ready
echo "Waiting for MySQL..."
while ! nc -z mysql 3306; do
  echo "MySQL not ready - waiting..."
  sleep 2
done
echo "MySQL is ready!"

cd /gocms/migrations
GOOSE_DRIVER="mysql" GOOSE_DBSTRING="root:root@tcp(mysql:3306)/gocms" goose up
cd /gocms
# go test ./... -v
air -c .air.toml