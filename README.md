# Exorcist

Similar to ghost and poltergeist this one is written in golang

## Getting started

- Install Go
- Install Docker
- Install psql
- Copy `templte.env` -> `.env` and fill in the details
- `docker compose up -d` to start the database
- `psql -U exorcist -h 127.0.0.1 -p 5432 -d exorcist -f ./migration/database.sql` initial database structure
- `make run` to start the application

## Tools

### Api

[Gin](https://go.dev/doc/tutorial/web-service-gin)

### Database

[raw postgres](https://golangdocs.com/golang-postgresql-example)

[Jet](https://github.com/go-jet/jet)

`jet -dsn=postgresql://${user}:${pass}@localhost:5432/exorcist?sslmode=disable -schema=public -path=./gen`
For some reason this does not run in zsh. Run it in bash

### Environment

[dotenv](https://github.com/joho/godotenv)