XC_OS="linux darwin"
XC_ARCH="amd64 arm64"
XC_PARALLEL="2"
BIN="dist"
SRC=$(shell find . -name "*.go")

ifeq (, $(shell which gox))
$(warning "could not find gox in $(PATH), run: go get github.com/mitchellh/gox")
endif

.PHONY: all build fmt test install_deps clean

default: all

all: fmt test build

build: install_deps
	gox \
		-os=$(XC_OS) \
		-arch=$(XC_ARCH) \
		-parallel=$(XC_PARALLEL) \
		-output=$(BIN)/{{.Dir}}_{{.OS}}_{{.Arch}} \
		;

fmt:
	$(info ******************** checking formatting ********************)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)
	terraform fmt -recursive


test: install_deps
	$(info ******************** running tests ********************)
	go test -v ./...

install_deps:
	$(info ******************** downloading dependencies ********************)
	go get -v ./...

clean:
	rm -rf $(BIN)
