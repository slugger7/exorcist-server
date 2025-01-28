#!/bin/bash
echo "Running migrations"

set -a && source ./.env && set +a

echo "Running migrations"
migrate -source file://./migrations  -database "${DATABASE_CONNECTION_STRING}" -verbose up
