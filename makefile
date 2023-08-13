default: air

build:
	go build -race -o bin/main src/main.go

test:
	go test ./...

air:
	air -build.cmd "make build" -build.bin "bin/main"
