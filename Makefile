NAME=gwvault

VERSION=$$(git describe --tags --always)

LDFLAGS=-ldflags=all="-X main.version=${VERSION}"

all: tools build

tools:
	go get -u -v "github.com/mitchellh/gox"

build:
	@mkdir -p bin/
	go get -t ./...
	go test -v ./...
	go build ${LDFLAGS} -o bin/${NAME} cmd/gwvault/main.go

xbuild: clean
	@mkdir -p build
	gox \
		-os="linux" \
		-os="windows" \
		-os="darwin" \
		-arch="amd64" \
		${LDFLAGS} \
		-output="build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)" \
		./...

package: xbuild
	$(eval FILES := $(shell ls build))
	@mkdir -p build/tgz
	for f in $(FILES); do \
		(cd $(shell pwd)/build && tar -zcvf tgz/$$f.tar.gz $$f); \
		echo $$f; \
	done

clean:
	@rm -rf bin/ && rm -rf build/

ci: tools package

.PHONY: all tools build xbuild package clean ci
