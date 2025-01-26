#!/bin/bash
docker compose down;
docker compose up -d;

echo "Waiting for database to spin up"
sleep 5;

PGPASSWORD=some-secure-password psql -U exorcist -h 127.0.0.1 -p 5432 -d exorcist -f ./migration/database.sql;
