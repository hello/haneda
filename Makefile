linux:
	GOOS=linux GOARCH=amd64 go build -o _build/server_linux_amd64
osx:
	GOOS=darwin GOARCH=amd64 go build -o _build/server_darwin_amd64

release: osx linux

.PHONY: release osx linux
