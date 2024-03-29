SOURCES := $(filter-out dockerfiles/%, $(wildcard */*))

BOLD         := $(shell tput -Txterm bold)
RED          := $(shell tput -Txterm setaf 1)
GREEN        := $(shell tput -Txterm setaf 2)
YELLOW       := $(shell tput -Txterm setaf 3)
BLUE         := $(shell tput -Txterm setaf 6)
NC := $(shell tput -Txterm sgr0)

BASEIMAGE = clinic
WORKSPACEIMAGE = hospital
IMAGESUFFIX = ward

all: base workspace
	@mkdir ~/projects

base: dockerfiles/base $(SOURCES)
	@docker build -t "$(BASEIMAGE)" -f dockerfiles/base

workspace: dockerfiles/workspace $(SOURCES)
	@docker build -t "$(WORKSPACEIMAGE)" -f dockerfiles/workspace

python-img:
	@docker build -t "python-$(IMAGESUFFIX)" -f dockerfiles/python

ts-img:
	@docker build -t "ts-$(IMAGESUFFIX)" -f dockerfiles/python

flutter-img:
	@docker build -t "flutter-$(IMAGESUFFIX)" -f dockerfiles/python

clean:
	@docker images -a | grep -q "$(BASEIMAGE)" && docker rmi "$(BASEIMAGE)"
	@docker images -a | grep -q "$(WORKSPACEIMAGE)" && docker rmi "$(WORKSPACEIMAGE)"
	@docker images -a | grep "-$(IMAGESUFFIX)$" | xargs docker rmi
