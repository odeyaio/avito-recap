#!/bin/sh
set -eu

mkdir -p /tmp/migrations
for f in /migrations/*.sql; do
  base=$(basename "$f")
  cp "$f" "/tmp/migrations/${base%.sql}.up.sql"
done

migrate -path=/tmp/migrations -database="postgres://${POSTGRES_USER:-avito-recap}:${POSTGRES_PASSWORD:-avito-recap}@postgres:5432/${POSTGRES_DB:-avito-recap}?sslmode=disable" up
