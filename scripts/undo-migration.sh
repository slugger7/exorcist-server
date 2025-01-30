#!/usr/bin/env bash
echo "Undoing migration"

./scripts/set-env.sh

echo "Running migrations"
migrate -source file://./migrations  -database "${DATABASE_CONNECTION_STRING}" -verbose down 1
