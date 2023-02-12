.PHONY: build
build:
	go build -o bin/terrascope cmd/terrascope/main.go

.PHONY: install
install:
	go build -o $$GOPATH/bin/terrascope cmd/terrascope/main.go
	# GOOS=linux GOARCH=amd64 go build -o ../../bin/terrascope_linux_amd64 cmd/terrascope/main.go
	# GOOS=linux GOARCH=arm64 go build -o ../../bin/terrascope_linux_arm64 cmd/terrascope/main.go
	# GOOS=darwin GOARCH=amd64 go build -o ../../bin/terrascope_darwin_amd64 cmd/terrascope/main.go
	# GOOS=darwin GOARCH=arm64 go build -o ../../bin/terrascope_darwin_arm64 cmd/terrascope/main.go

.PHONY: test
test:
	go test ./...
