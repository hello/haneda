linux:
	GOOS=linux GOARCH=amd64 go build -o _build/haneda_osx_amd64
osx:
	GOOS=darwin GOARCH=amd64 go build -o _build/haneda_osx_amd64

release: osx linux

.PHONY: release osx linux
