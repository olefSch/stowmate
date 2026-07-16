VERSION := `git describe --tags --always --dirty 2>/dev/null || echo dev`
COMMIT := `git rev-parse --short HEAD 2>/dev/null || echo unknown`
DATE := `date -u +%Y-%m-%dT%H:%M:%SZ`
LDFLAGS := '-s -w'

build:
    go build -trimpath -buildvcs=false -ldflags="{{LDFLAGS}} -X main.version={{VERSION}} -X main.commit={{COMMIT}} -X main.date={{DATE}}" -o stowmate ./cmd/stowmate/

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

docs-serve:
    uvx zensical serve

docs-build:
    uvx zensical build

default:
    @just --list
