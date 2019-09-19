TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build

build: all

fmt:
	@gofmt -w $(GOFMT_FILES)
	@goimports -w $(GOFMT_FILES)

vendor:
	@govendor sync

.PHONY: build fmt vendor

all: fmt vendor mac windows linux

dev: clean fmt mac copy

copy:
	tar -xvf bin/BlackHole_darwin-amd64.tgz

clean:
	rm -rf bin/*

mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/BlackHole
	tar czvf bin/BlackHole_darwin-amd64.tgz bin/BlackHole
	rm -rf bin/BlackHole

windows:
	GOOS=windows GOARCH=amd64 go build -o bin/BlackHole.exe
	tar czvf bin/BlackHole_windows-amd64.tgz bin/BlackHole.exe
	rm -rf bin/BlackHole.exe

linux:
	GOOS=linux GOARCH=amd64 go build -o bin/BlackHole
	tar czvf bin/BlackHole_linux-amd64.tgz bin/BlackHole
	rm -rf bin/BlackHole
