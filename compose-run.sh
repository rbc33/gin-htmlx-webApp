#! /usr/bin/env bash
apk add git
set -euo pipefail
git config --global --add safe.directory '*'
cd /gocms/migrations
GOOSE_DRIVER="mysql" GOOSE_DBSTRING="root:root@tcp(mysql:3306)/gocms" goose up
cd /gocms
air
