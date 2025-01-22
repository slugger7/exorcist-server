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
- `psql -U exorcist -h 127.0.0.1 -p 5432 -d exorcist -f ./migration/database.sql` initial database structure
- `make run` to start the application

## Tools

### Api

[Gin](https://go.dev/doc/tutorial/web-service-gin)

### Database

For this application I want to try creating a database first design while also not using an ORM but rather a SQL builder.

There is still some work to do to figure out how to properly create migrations for this project though.

[raw postgres](https://golangdocs.com/golang-postgresql-example)

[Jet](https://github.com/go-jet/jet)

`jet -dsn=postgresql://${user}:${pass}@localhost:5432/exorcist?sslmode=disable -schema=public -path=./gen`
For some reason this does not run in zsh. Run it in bash

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
