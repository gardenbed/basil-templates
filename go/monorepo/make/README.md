# Make

## Documentation

| File | Required Variables | Macros | Rules |
|------|--------------------|--------|-------|
| `common.mk` | | `echo_red` <br/> `echo_green` <br/> `echo_yellow` <br/> `echo_blue` <br/> `echo_purple` <br/> `echo_cyan` | |
| `go.mk` | `name` | | `test` <br/> `test-short` <br/> `test-coverage` <br/> `clean-test` <br/> `run` <br/> `build` <br/> `build-all` <br/> `clean-build` |
| `grpc.mk` | `proto_path` <br/> `go_out_path` | | `protoc` <br/> `protoc-gen-go` <br/> `protobuf` |
| `docker.mk` | `docker_image` <br/> `docker_tag` | | `docker` <br/> `docker-test` <br/> `push` <br/> `push-latest` <br/> `save-docker` <br/> `load-docker` <br/> `clean-docker` |
| `terraform.mk` | | `create_aws_key` <br/> `create_gcp_key` | `validate` <br/> `plan` <br/> `apply` <br/> `refresh` <br/> `destroy` <br/> `clean-terraform` |
