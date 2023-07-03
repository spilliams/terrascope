version := 1.0.0
githash := $(shell git rev-parse --short HEAD)
buildtime := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p_%Z')
ldflags := "\
	-X github.com/spilliams/terrascope/internal/version.versionNumber=$(version)\
	-X github.com/spilliams/terrascope/internal/version.gitHash=$(githash)\
	-X github.com/spilliams/terrascope/internal/version.buildTime=$(buildtime)"

.PHONY: build
build:
	go build -ldflags $(ldflags) -o bin/terrascope main.go

.PHONY: install
install:
	go build -ldflags "${ldflags}" -o $$GOPATH/bin/terrascope main.go

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	# binary will be $(go env GOPATH)/bin/golangci-lint
	# curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.53.3

	golangci-lint run -v
