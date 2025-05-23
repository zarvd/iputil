BIN ?= $(shell pwd)/bin
GOLANGCI_LINT_VERSION ?= v2.0.2
GOLANGCI_LINT = $(BIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)/golangci-lint
GOLANGCI_LINT_CONFIG ?= $(shell pwd)/.golangci.yaml

all: help

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

lint: golangci-lint ## Run the linter.
	$(GOLANGCI_LINT) run --config $(GOLANGCI_LINT_CONFIG)

test: ## Run the tests.
	go test ./...

##@ Toolings

golangci-lint: ## Install the linter.
	$(call go-get-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION))

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
	set -e ;\
	TMP_DIR=$$(mktemp -d) ;\
	cd $$TMP_DIR ;\
	go mod init tmp ;\
	echo "Downloading $(2)" ;\
	GOBIN=$(shell dirname $(1)) go install $(2) ;\
	rm -rf $$TMP_DIR ;\
}
endef
