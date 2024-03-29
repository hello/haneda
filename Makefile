PACKAGES = $(shell go list ./... | grep -v '/vendor/')
linux:
	GOOS=linux GOARCH=amd64 go build -o _build/haneda_linux_amd64
osx:
	GOOS=darwin GOARCH=amd64 go build -o _build/haneda_osx_amd64

format:
	@echo "--> Running go fmt"
	@go fmt $(PACKAGES)

test:
	go test ./core/
	go test ./sense/
	go test ./client/
	go test .

release: osx linux

.PHONY: release osx linux
