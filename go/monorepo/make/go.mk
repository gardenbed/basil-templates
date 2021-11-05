## PREREQUISITES:
##
## The following variables need to be defined where this file is included:
##   - name
##   - main_pkg
##
## The following files need to be included where this file is included:
##   - common.mk
##


## MACROS
##

# Compiles a binary for the host platform.
define compile
  go build -ldflags $(ldflags) -o $(build_dir)/$(name) $(main_pkg)
	$(call echo_green,$(build_dir)/$(name));
endef

# Compiles a binary for a target platform.
define cross_compile
	GOOS=$(shell echo $(1) | cut -d '-' -f 1) \
	GOARCH=$(shell echo $(1) | cut -d '-' -f 2) \
	go build -ldflags $(ldflags) -o $(build_dir)/$(name)-$(1) $(main_pkg)
	$(call echo_green,$(build_dir)/$(name)-$(1));
endef


## VARIABLES
##

make_dir := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

build_dir := bin
platforms := linux-386 linux-amd64 linux-arm linux-arm64 darwin-amd64 windows-386 windows-amd64


## BUILD METADATA & FLAGS
##
## Git tags are leveraged for semantic versioning.
## A package named version is expected to exist with the following file inside it:
##
## package version
##
## var (
##   Version   string
##   Commit    string
##   Branch    string
##   GoVersion string
##   BuildTool string
##   BuildTime string
## )
##

version := $(shell $(make_dir)/semver.sh)
commit := $(shell git rev-parse --short HEAD)
branch := $(shell git rev-parse --abbrev-ref HEAD)
go_version := $(shell go version | grep -E -o '[0-9]+\.[0-9]+\.[0-9]+')
build_tool := Make
build_time := $(shell date '+%Y-%m-%d %T %Z')

metadata_package := $(shell go list ./... | grep -E 'metadata$$' | head -n 1)
version_flag := -X "$(metadata_package).Version=$(version)"
commit_flag := -X "$(metadata_package).Commit=$(commit)"
branch_flag := -X "$(metadata_package).Branch=$(branch)"
go_version_flag := -X "$(metadata_package).GoVersion=$(go_version)"
build_tool_flag := -X "$(metadata_package).BuildTool=$(build_tool)"
build_time_flag := -X "$(metadata_package).BuildTime=$(build_time)"
ldflags := '$(version_flag) $(commit_flag) $(branch_flag) $(go_version_flag) $(build_tool_flag) $(build_time_flag)'


## RULES
##

.PHONY: test
test:
	go test -race ./...

.PHONY: test-short
test-short:
	go test -short ./...

.PHONY: test-coverage
test-coverage:
	go test -covermode=atomic -coverprofile=c.out ./...
	go tool cover -html=c.out -o coverage.html

.PHONY: clean-test
clean-test:
	rm -f c.out coverage.html

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	@ $(call compile)

.PHONY: build-all
build-all:
	@ $(foreach platform, $(platforms), $(call cross_compile,$(platform)))

.PHONY: clean-build
clean-build:
	rm -rf $(name) $(build_dir)
