.PHONY: main
main: clean test lint release

PKGS := $(shell go list ./...)
.PHONY: test
test:
	go test $(PKGS)

BIN_DIR := $(GOPATH)/bin
GOLINTER := $(BIN_DIR)/golint

golint := $(if $(GOPATH),$(GOPATH),$(HOME)/go)/bin/golint
$(GOLINTER):
	go get -u golang.org/x/lint/golint

.PHONY: lint
lint: $(GOLINTER)
	$(golint) -set_exit_status ./...

clean:
	rm -rf ./release/

VERSION := $(shell go run cmd/dddd.go -m)
.PHONY: version
version:
	go run cmd/dddd.go -m

BINARY := dddd
PLATFORMS := windows linux darwin
os = $(word 1, $@)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	GOOS=$(os) GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-$(os)-amd64 cmd/dddd.go

.PHONY: micro
micro:
	mkdir -p release
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags "$(govvv -flags)" -o release/microdddd cmd/dddd.go 

CONTAINER_BUILDER := $(shell which docker && printf docker || printf podman)
SUDO_REQUIRED := sudo

.PHONY: container
container:
	$(SUDO_REQUIRED) $(CONTAINER_BUILDER) build -f packaging/docker/Dockerfile -t tcarrio/dddd:$(VERSION) .

.PHONY: release
release: windows linux darwin


