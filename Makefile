# Copyright 2019 AppsCode Inc.

SHELL=/bin/bash -o pipefail

GO_PKG   := appscode-cloud
REPO     := $(notdir $(shell pwd))

SRC_PKGS := api cmd
SRC_DIRS := $(SRC_PKGS) data *.go # directories which hold app source (not vendored)

# Used internally.  Users should pass GOOS and/or GOARCH.
OS   := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))

GO_VERSION       ?= 1.25
BUILD_IMAGE      ?= ghcr.io/appscode/golang-dev:$(GO_VERSION)

fmt: $(BUILD_DIRS)
	@docker run                                                 \
		-i                                                      \
		--rm                                                    \
		-u $$(id -u):$$(id -g)                                  \
		-v $$(pwd):/src                                         \
		-w /src                                                 \
		-v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin                \
		-v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH)  \
		-v $$(pwd)/.go/cache:/.cache                            \
		--env HTTP_PROXY=$(HTTP_PROXY)                          \
		--env HTTPS_PROXY=$(HTTPS_PROXY)                        \
		$(BUILD_IMAGE)                                          \
		./hack/fmt.sh $(SRC_DIRS)

.PHONY: gen
gen:
	@true

publish:
	@echo "publishing files"
	aws s3 sync $$(pwd)/files s3://cdn-appscode-com/files --delete --acl public-read
	@echo "publishing images"
	aws s3 sync $$(pwd)/images s3://cdn-appscode-com/images --delete --acl public-read
