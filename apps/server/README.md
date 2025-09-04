# Exorcist

Similar to [ghost](https://github.com/slugger7/ghost-media) and [poltergeist](https://github.com/slugger7/poltergeist) this one is written in golang

## Under construction

This is your warning that this project is not in any state ready to be used (even by me) and is still under constant construction without a usable version at the moment

## Getting started

- Install Go
- Install Docker
- Install psql
  - Mac:
    - `brew install libpq`
    - `brew link --force libpq`
- Install ffmpeg
- Install [tygo](https://github.com/gzuidhof/tygo)
- Copy `.env.example` -> `.env` and fill in the details
- `docker compose up -d` to start the database
- `make run` to start the application

## Frontend

The server will serve any files that are in the [www](./www) directory if it exists. This directory can be changed by an environment variable in [.env](.env) but it is an optional field.
If the variable is not set it will not serve any static files.

The [web](https://github.com/slugger7/exorcist-web) project can be built and all of the files in its `dist` directory can be added to the [www](./www) directory.

One thing to keep in mind is that any path that does not have the `/api` prefix will be forwarded to the webserver as it is not something for this webserver to handle and it assumes that it is meant for the frontend portion of the project like a Single Page Application

## Tools

### Database

#### Entity Relation Diagram

![entity_relation_diagram](./diagrams/out/entity_relation_diagram.d2.svg)

This application is using a database first approach with a sql builder instead of an ORM.
It utilizes [Jet](https://github.com/go-jet/jet)

#### Migrations

We utilize [golang-migrate](https://github.com/golang-migrate/migrate) for our migrations

Migrations run automatically when the application starts up.
It is recommended to install the cli tool for running migrations. This will allow you to run migrations from the command line without having to start up the application. [CLI migration](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
It is also recommended that you install the [Jet CLI tool](https://github.com/go-jet/jet?tab=readme-ov-file#prerequisites) in order to update the models of the application

- Create migration: `./scripts/create-migration.sh <migration-name>`
- Run migrations: `./scripts/run-migrations.sh`
- Undo a migration: `./scripts/undo-migration.sh`
- Update models: `./scripts/update-models.sh`

The usual workflow would to add a migration would be:

1. Create a migration
1. Run the migrations
1. Update the models

### Troubleshooting

- If using zsh remember to add the following to your .zshrc

  ```bash
  export GOPATH=$HOME/go  
  export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
  ```

### Diagrams

Generating of diagrams is done by utilizisg [d2](https://d2lang.com/)
In order to generate diagrams you will need to install the cli tool

### Mocks

We are using a mock generating library to prevent us having to create these mocks by hand. To get this working:

- Install the cli tool [mockgen](https://github.com/uber-go/mock)
- When an interface has changed you should be able to run `make mocks` and the mocks will be regenerated for you


## Troubleshooting

### zsh

Add the following to your `.zshrc`

```bash
GOPATH=$HOME/go  PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```

To activate it on the current shell:

```bash
source .zshrc
```

