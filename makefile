run:
	go run ./cmd/exorcist

get:
	go get ./cmd/exorcist

update:
	go get -u ./...

test:
	@go test ./... -v

exorcist:
	go build -o ./build/exorcist ./cmd/exorcist

build: exorcist

clean:
	rm -rf build
	rm -rf ts/*

mocks:
	./scripts/generate-mocks.sh

run-migrations:
	./scripts/run-migrations.sh

undo-migration:
	./scripts/undo-migration.sh

update-models:
	./scripts/update-models.sh

recreate-db:
	docker compose down db
	docker compose up -d
	
generate-diagrams:
	./scripts/generate-diagrams.sh

dtos:
	tygo generate
