TERRAFORM ?= terraform
COMPOSE ?= docker-compose

all: fmt syntax docs

docs:
	terraform-docs markdown . > README.md


$(TERRAFORM) fmt -recursive

syntax:
	$(TERRAFORM) fmt -diff -check -recursive

tests-compose:
	COMPOSE=$(COMPOSE) ./tests/main.sh local
tests-ci:
	./tests/main.sh ci


