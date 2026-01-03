.PHONY: build test clean run deps

build:
	go build -o bin/codigosH ./cmd/codigosH

test:
	go test ./...

clean:
	rm -rf bin/

run: build
	./bin/codigosH

deps:
	go mod download
	@echo "Dependencias instaladas"
	go mod tidy
	go mod download