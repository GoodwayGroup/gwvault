NAME=gwvault

VERSION=$$(git describe --tags --always)
SHORT_VERSION=$$(git describe --tags --always | awk -F '-' '{print $$1}')

LDFLAGS=-ldflags=all="-X main.version=${SHORT_VERSION}"

all: tools build

tools:
	GO111MODULE=off go get -u -v "github.com/mitchellh/gox"

build:
	@mkdir -p bin/
	go get -t ./...
	go test -v ./...
	go build ${LDFLAGS} -o bin/${NAME} ./main.go

xbuild: clean
	@mkdir -p build
	gox \
		-os="linux" \
		-os="windows" \
		-os="darwin" \
		-arch="amd64" \
		${LDFLAGS} \
		-output="build/$(NAME)_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)" \
		./...

package: xbuild
	$(eval FILES := $(shell ls build))
	@mkdir -p build/tgz
	for f in $(FILES); do \
		(cd $(shell pwd)/build && tar -zcvf tgz/$$f.tar.gz $$f); \
		echo $$f; \
	done

docs:
	DOCS_MD=1 go run ./main.go > docs/${NAME}.md
	DOCS_MAN=1 go run ./main.go > docs/${NAME}.8

clean:
	@rm -rf bin/ && rm -rf build/

ci: tools package

.PHONY: all tools build xbuild package clean ci docs
