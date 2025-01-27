run:
	go run ./cmd/exorcist

get:
	go get ./cmd/exorcist

test:
	@go test ./... -v

exorcist:
	go build ./cmd/exorcist

build: exorcist

clean:
	rm exorcist
