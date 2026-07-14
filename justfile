build:
    go build -o stowmate ./cmd/stowmate/

test:
    go test -race -cover ./...

test-e2e:
    go test -tags e2e -race ./...

lint:
    golangci-lint run ./...

vet:
    go vet ./...

clean:
    rm -f stowmate
    rm -rf dist/

install: build
    cp stowmate $(go env GOPATH)/bin/

coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

tidy:
    go mod tidy

vuln:
    go run golang.org/x/vuln/cmd/govulncheck@latest ./...

check: vet lint test vuln

default:
    @just --list
