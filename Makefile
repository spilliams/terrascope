version := 0.1.1
githash := $(shell git rev-parse --short HEAD)
buildtime := $(shell date -u '+%Y-%m-%d_%I:%M:%S%pm_%Z')
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
