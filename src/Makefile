.PHONY: build
build:
	go build -o bin/terraboots cmd/terraboots/main.go

.PHONY: install
install:
	go build -o $$GOPATH/bin/terraboots cmd/terraboots/main.go
	# GOOS=linux GOARCH=amd64 go build -o ../../bin/terraboots_linux_amd64 cmd/terraboots/main.go
	# GOOS=linux GOARCH=arm64 go build -o ../../bin/terraboots_linux_arm64 cmd/terraboots/main.go
	# GOOS=darwin GOARCH=amd64 go build -o ../../bin/terraboots_darwin_amd64 cmd/terraboots/main.go
	# GOOS=darwin GOARCH=arm64 go build -o ../../bin/terraboots_darwin_arm64 cmd/terraboots/main.go

.PHONY: test
test:
	go test ./...
