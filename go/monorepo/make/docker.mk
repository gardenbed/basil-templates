## PREREQUISITES:
##
## The following variables need to be defined where this file is included:
##   - docker_image
##   - docker_tag
##


## RULES
##

docker:
	@ docker image build \
	    --build-arg ldflags=$(ldflags) \
	    --tag $(docker_image):$(docker_tag) \
	    .

.PHONY: docker-test
docker-test:
	@ docker image build \
	    --file Dockerfile.test \
	    --tag $(docker_image):$(docker_tag) \
	    .

.PHONY: push
push:
	docker image push $(docker_image):$(docker_tag)

.PHONY: push-latest
push-latest:
	docker image tag $(docker_image):$(docker_tag) $(docker_image):latest
	docker image push $(docker_image):latest

.PHONY: save-docker
save-docker:
	docker image save -o docker.tar $(docker_image):$(docker_tag)

.PHONY: load-docker
load-docker:
	docker image load -i docker.tar

.PHONY: clean-docker
clean-docker:
	rm -f docker.tar
