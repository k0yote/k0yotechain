build:
	go build -o ./bin/k0yote

run: build
	./bin/k0yote

test:
	go test ./...