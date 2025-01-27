NAME=image-builder
BINARY=packer-plugin-${NAME}
# https://github.com/hashicorp/packer-plugin-sdk/issues/187
HASHICORP_PACKER_PLUGIN_SDK_VERSION?="v0.5.2"
PLUGIN_FQN=$(shell grep -E '^module' <go.mod | sed -E 's/module \s*//')
PLUGIN_PATH=./cmd/plugin

.PHONY: build
build:
	@go build -ldflags="-X '${PLUGIN_FQN}/main.Version=$(shell git describe --tags --abbrev=0)'" -o ${BINARY} ${PLUGIN_PATH}

.PHONY: dev
dev: build
	packer plugins install --path ${BINARY} "$(shell echo "${PLUGIN_FQN}" | sed 's/packer-plugin-//')"

.PHONY: test
test:
	@go test -race ./...

.PHONY: install-packer-sdc
install-packer-sdc:
	@go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@${HASHICORP_PACKER_PLUGIN_SDK_VERSION}

.PHONY: plugin-check
plugin-check: install-packer-sdc build
	@packer-sdc plugin-check ${BINARY}

.PHONY: generate
generate: install-packer-sdc
# https://github.com/hashicorp/packer-plugin-sdk/issues/187
	@go mod edit -replace "github.com/zclconf/go-cty=github.com/nywilken/go-cty@v1.13.3"
	@go mod tidy
	@go generate ${PLUGIN_PATH}
