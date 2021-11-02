# Include macros, variables, and rules
include ../monorepo/make/common.mk
include ../monorepo/make/go.mk      # test, test-short, test-coverage, clean-test, run, build, build-all, clean-build
include ../monorepo/make/docker.mk  # docker, docker-test, push, push-latest, save-docker, load-docker, clean-docker

# Variables required by inclusions
name := http-service-vertical
docker_image := dockerid/http-service-vertical
docker_tag ?= $(version)
