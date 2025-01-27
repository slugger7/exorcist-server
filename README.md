# Exorcist

Similar to ghost and poltergeist this one is written in golang

## Getting started

- Install Go
- Install Docker
- Install psql
  - Mac:
    - `brew install libpq`
    - `brew link --force libpq`
- Install ffmpeg
- Copy `templte.env` -> `.env` and fill in the details
- `docker compose up -d` to start the database
- `make run` to start the application

## Tools

### Api

[Gin](https://go.dev/doc/tutorial/web-service-gin)

### Database

For this application I want to try creating a database first design while also not using an ORM but rather a SQL builder.

There is still some work to do to figure out how to properly create migrations for this project though.

[raw postgres](https://golangdocs.com/golang-postgresql-example)

[Jet](https://github.com/go-jet/jet)

`jet -source=postgres -host=localhost -port=5432 -user=exorcist -password=some-secure-password -dbname=exorcist -schema=public -sslmode=disable -path=./internal/db`
If you use zsh look at the troubleshooting section for adding gopath to your path

#### Migrations

For migrations look into [golang-migrate](https://github.com/golang-migrate/migrate)

You should never have to manually run migrations as it should run when the application starts up.
To make life easier creating migrations you can install the cli of the migration tool from [here](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

- Create migration: `migrate create -ext=sql -dir=./migrations <migration-name>`
- Run migrations: `./run-migrations.sh`

### FFMpeg stuff

[Wrapper](https://github.com/u2takey/ffmpeg-go)

### Environment

[dotenv](https://github.com/joho/godotenv)

### Troubleshooting

- If using zsh remember to add the following to your .zshrc

  ```bash
  export GOPATH=$HOME/go  
  export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
  ```
