GO_MOD_NAME := $(shell grep -m 1 'module ' go.mod | cut -d' ' -f2)
GO_MOD_DOMAIN := $(shell echo $(GO_MOD_NAME) | awk -F '/' '{print $$1}')
PROJECT_NAME := $(shell grep 'module ' go.mod | awk '{print $$2}' | sed 's|$(GO_MOD_DOMAIN)/||g')
NOW := $(shell date +'%Y%m%d%H%M%S')

TCR_HOST_LOCAL := costa92

TCR_HOST ?= $(TCR_HOST_LOCAL)

TCR_IMAGE_LOCAL := $(TCR_HOST)/$(shell echo $(PROJECT_NAME) | sed 's^\/^\/local\/^')

 LOCAL_TAG := "local-"$(NOW)
TAG := $(shell git describe --always --tags --abbrev=0 --match 'v*' --exclude '*/*' | tr -d "[\r\n]")
.PHONY: build
#
local/build:
	docker build -q -t $(TCR_IMAGE_LOCAL):$(LOCAL_TAG) .
#
prod/build:
	docker build -q -t $(TCR_IMAGE_LOCAL):$(TAG) .
	docker push $(TCR_IMAGE_LOCAL):$(TAG)


fmt:
	@gofumpt -version || go install mvdan.cc/gofumpt@latest
	gofumpt -extra -w -d .
	@gci -v || go install github.com/daixiang0/gci@latest
	gci write -s standard -s default -s 'Prefix($(GO_MOD_DOMAIN))' --skip-generated .