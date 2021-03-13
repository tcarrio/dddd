.PHONY: main
main: clean test lint build

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
	rm -rf $(RELEASES) $(ARTIFACTS)

VERSION := $(shell go run cmd/dddd.go -m)
.PHONY: version
version:
	go run cmd/dddd.go -m

ARTIFACTS := ./artifacts
RELEASES := ./release

CHANGELOG := changelog.md
.PHONY: changelog
changelog:
	mkdir -p $(ARTIFACTS)
	echo "dddd v$(VERSION)" > $(ARTIFACTS)/$(CHANGELOG)

BINARY := dddd
PLATFORMS := windows linux darwin
os = $(word 1, $@)

.PHONY: $(PLATFORMS)
$(PLATFORMS): prebuild
	GOOS=$(os) GOARCH=amd64 go build -o $(RELEASES)/$(BINARY)-$(VERSION)-$(os)-amd64 cmd/dddd.go

.PHONY: micro
micro: prebuild
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags "$(govvv -flags)" -o $(RELEASES)/microdddd cmd/dddd.go 

CONTAINER_BUILDER := $(shell which docker && printf docker || printf podman)
SUDO_REQUIRED := sudo

.PHONY: container
container:
	$(SUDO_REQUIRED) $(CONTAINER_BUILDER) build -f packaging/docker/Dockerfile -t tcarrio/dddd:$(VERSION) .

.PHONY: prebuild
prebuild:
	mkdir -p $(RELEASES)

.PHONY: build
build: prebuild windows linux darwin

DRAFT_MODE := $(if $(RELEASE_PROD),,-d)
.PHONY: release
release: hub clean build changelog
	gh release view $(VERSION) || \
	gh release create v$(VERSION) $(DRAFT_MODE) -F $(ARTIFACTS)/$(CHANGELOG) $(RELEASES)/*

.PHONY: hub
hub:
	which gh || which hub