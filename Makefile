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

BINARY := dddd
VERSION := $(shell go run cmd/dddd.go -m)
PLATFORMS := windows linux darwin
os = $(word 1, $@)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	GOOS=$(os) GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-$(os)-amd64 cmd/dddd.go 

.PHONY: release
release: windows linux darwin


