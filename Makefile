NAME=gwvault

VERSION=$$(git describe --tags --always)
SHORT_VERSION=$$(git describe --tags --always | awk -F '-' '{print $$1}')

LDFLAGS=-ldflags=all="-X main.version=${SHORT_VERSION}"

all: build

build:
	@mkdir -p bin/
	go get -t ./...
	go test -v ./...
	go build ${LDFLAGS} -o bin/${NAME} ./main.go

docs:
	DOCS_MD=1 go run ./main.go > docs/${NAME}.md
	DOCS_MAN=1 go run ./main.go > docs/${NAME}.8

clean:
	@rm -rf bin/ && rm -rf build/ && rm -rf dist/

.PHONY: all tools build clean docs
