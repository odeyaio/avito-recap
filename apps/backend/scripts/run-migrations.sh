#!/bin/sh
set -eu

migrate -path=/migrations -database="postgres://${POSTGRES_USER:-avito-recap}:${POSTGRES_PASSWORD:-avito-recap}@postgres:5432/${POSTGRES_DB:-avito-recap}?sslmode=disable" "${@:-up}"
