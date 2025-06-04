#!/usr/bin/env bash
echo "Running migrations"

. ./scripts/set-env.sh

echo "Running migrations"
migrate -path ./migrations  -database "${DATABASE_CONNECTION_STRING}" -verbose up
