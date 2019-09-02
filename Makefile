# Copyright 2019 AppsCode Inc.

SHELL=/bin/bash -o pipefail

GO_PKG   := appscode-cloud
REPO     := $(notdir $(shell pwd))

SRC_PKGS := api cmd
SRC_DIRS := $(SRC_PKGS) data # directories which hold app source (not vendored)

# Used internally.  Users should pass GOOS and/or GOARCH.
OS   := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))

GO_VERSION       ?= 1.12.9
BUILD_IMAGE      ?= appscode/golang-dev:$(GO_VERSION)-stretch

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
gen: gen-bindata

.PHONY: gen-bindata
gen-bindata:
	@docker run                                                 \
		-i                                                      \
		--rm                                                    \
		-u $$(id -u):$$(id -g)                                  \
		-v $$(pwd):/src                                         \
		-w /src                                                 \
		-v /tmp:/.cache                                         \
		--env HTTP_PROXY=$(HTTP_PROXY)                          \
		--env HTTPS_PROXY=$(HTTPS_PROXY)                        \
		$(BUILD_IMAGE)                                          \
		/bin/bash -c "                                          \
			cd /src/data;                                                                                \
			go-bindata -ignore=\\.go -ignore=\\.DS_Store -mode=0644 -o bindata.go -pkg data ./...;       \
			cd /src/data/products;                                                                       \
			go-bindata -ignore=\\.go -ignore=\\.DS_Store -mode=0644 -o bindata.go -pkg products ./...;   \
			cd /src/hugo;                                                                                \
			go-bindata -ignore=\\.go -ignore=\\.DS_Store -mode=0644 -o bindata.go -pkg hugo ./...;       \
		"

publish:
	@echo "publishing files"
	gsutil rsync -d -r $$(pwd)/files gs://appscode-cdn/files
	gsutil acl ch -u AllUsers:R -r gs://appscode-cdn/files
	@echo "publishing images"
	gsutil rsync -d -r $$(pwd)/images gs://appscode-cdn/images
	gsutil acl ch -u AllUsers:R -r gs://appscode-cdn/images
