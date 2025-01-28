#!/bin/bash
echo "Updating models"

set -a && source ./.env && set +a

jet -source=postgres -host=${DATABASE_HOST} -port=${DATABASE_PORT} -user=${DATABASE_USER} -password=${DATABASE_PASSWORD} -dbname=${DATABASE_NAME} -schema=public -sslmode=disable -path=./internal/db
