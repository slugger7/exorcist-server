#!/bin/bash
echo "Undoing migration"

set -a && source ./.env && set +a

echo "Running migrations"
migrate -source file://./migrations  -database "${DATABASE_CONNECTION_STRING}" -verbose down 1
