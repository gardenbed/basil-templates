## PREREQUISITES:
##


## MACROS
##

# Creates an SSH key pair for AWS
define create_aws_key
	@ mkdir -p $(shell dirname $(1))
	@ ssh-keygen -f $(1) -t rsa -N '' 1> /dev/null
	@ chmod 400 $(1)
	@ mv $(1) $(1).pem
endef

# Creates an SSH key pair for Google Cloud
define create_gcp_key
	@ mkdir -p $(shell dirname $(1))
	@ ssh-keygen -f $(1) -t rsa -N '' -C $(2) 1> /dev/null
	@ chmod 400 $(1)
	@ mv $(1) $(1).pem
endef


## VARIABLES
##

uuid := $(shell uuidgen | tr [:upper:] [:lower:])
owner := $(shell whoami)
branch := $(shell git rev-parse --abbrev-ref HEAD)
commit := $(shell git rev-parse --short HEAD)


## RULES
##

.PHONY: validate
validate:
	@ terraform validate

.PHONY: plan
plan:
	terraform plan \
	  -var uuid=$(uuid) \
	  -var owner=$(owner) \
	  -var git_branch=$(branch) \
	  -var git_commit=$(commit)

.PHONY: apply
apply:
	terraform apply \
	  -var uuid=$(uuid) \
	  -var owner=$(owner) \
	  -var git_branch=$(branch) \
	  -var git_commit=$(commit)

.PHONY: refresh
refresh:
	terraform refresh \
	  -var uuid=$(uuid) \
	  -var owner=$(owner) \
	  -var git_branch=$(branch) \
	  -var git_commit=$(commit)

.PHONY: destroy
destroy:
	terraform destroy \
	  -var uuid=$(uuid) \
	  -var owner=$(owner) \
	  -var git_branch=$(branch) \
	  -var git_commit=$(commit)

.PHONY: clean-terraform
clean-terraform:
	@ rm -rf .terraform terraform.tfstate terraform.tfstate.backup
