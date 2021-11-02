## PREREQUISITES:
##
## The following variables need to be defined where this file is included:
##   - proto_path
##   - go_out_path
##
## Proto files belonging to the same package expected to be in a directory with the same name as the package name.
##


## MACROS
##

# Determines the operating system and architecture
define get_os_arch
	arch = $(shell uname -m)
	os_name = $(shell uname -s | tr [:upper:] [:lower:])
	ifeq ($$(os_name),linux)
		$(1) = $$(os_name)-$$(arch)
	else ifeq ($$(os_name),darwin)
		$(1) = osx-$$(arch)
	endif
endef


## VARIABLES
##

$(eval $(call get_os_arch,os_arch))
protoc_release = $(shell curl -s https://api.github.com/repos/protocolbuffers/protobuf/releases/latest | jq -r '.tag_name')
protoc_version = $(shell echo $(protoc_release) | cut -c2-)


## RULES
##

.PHONY: check-tools
check-tools:
	@ command -v curl
	@ command -v jq
	@ command -v unzip
	@ command -v git
	@ command -v go

.PHONY: protoc
protoc: check-tools
	curl -fsSL https://github.com/protocolbuffers/protobuf/releases/download/$(protoc_release)/protoc-$(protoc_version)-$(os_arch).zip -o protoc.zip
	unzip -o protoc.zip -d /usr/local bin/protoc
	unzip -o protoc.zip -d /usr/local include/*
	rm -f protoc.zip

.PHONY: protoc-gen-go
protoc-gen-go:
	go install github.com/golang/protobuf/protoc-gen-go@latest

.PHONY: protobuf
protobuf:
	@ mkdir -p $(go_out_path)
	protoc \
	  --proto_path=$(proto_path) \
	  --go_out=paths=source_relative,plugins=grpc:$(go_out_path) \
	  $(foreach proto_file, $(shell find $(proto_path) -name '*.proto'), $(proto_file))
