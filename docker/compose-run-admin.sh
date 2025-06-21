#! /usr/bin/env bash
apk add git
set -euo pipefail
git config --global --add safe.directory '*'
cd /gocms
# go test ./... -v
air -c .air.admin.toml