linux:
	GOOS=linux GOARCH=amd64 go build -o _build/haneda_linux_amd64
osx:
	GOOS=darwin GOARCH=amd64 go build -o _build/haneda_osx_amd64

test:
	go test -v ./core/
	go test -v ./sense/
	go test -v .
release: osx linux

.PHONY: release osx linux
