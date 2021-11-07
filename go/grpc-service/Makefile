# Include macros, variables, and rules
include ../monorepo/make/common.mk
include ../monorepo/make/go.mk      # test, test-short, test-coverage, clean-test, run, build, build-all, clean-build
include ../monorepo/make/grpc.mk    # check-tools, protoc, protoc-gen-go, protobuf
include ../monorepo/make/docker.mk  # docker, docker-test, push, push-latest, save-docker, load-docker, clean-docker

# Variables required by inclusions
name := grpc-service
main_pkg := .
proto_path := idl
go_out_path := internal/idl
docker_image := dockerid/grpc-service
docker_tag ?= $(version)
