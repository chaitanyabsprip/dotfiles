.PHONY: clean all

BOLD         := $(shell tput -Txterm bold)
RED          := $(shell tput -Txterm setaf 1)
GREEN        := $(shell tput -Txterm setaf 2)
YELLOW       := $(shell tput -Txterm setaf 3)
BLUE         := $(shell tput -Txterm setaf 4)
NC           := $(shell tput -Txterm sgr0)

BASEIMAGE = clinic
WORKSPACEIMAGE = hospital
IMAGESUFFIX = ward
SOURCES := $(filter-out dockerfiles/%, $(wildcard */*))
SOURCES := $(filter-out $(wildcard $(shell find . -type l -printf '%P\n')), $(SOURCES))

define build_image
	@echo "$(BLUE)$(BOLD)Building $(1) image$(NC)"
	@$(if $(SKIP),docker image inspect "$(1)" >/dev/null 2>&1 \
		&& echo "$(YELLOW)Image already exists$(NC)" && exit;)\
		docker build $(if $(DEBUG),--progress=plain) -t "$(1)" -t "chaitanyabsprip/$(1)" -f $(2) .
endef

define clean_image
	@$(if $(SKIP),exit;)echo "$(RED)$(BOLD)Deleting $(1) image$(NC)" \
		&& docker images -a --format '{{.Repository}}' | grep -w "$(1)" \
		| xargs -I {} docker rmi {} >/dev/null 2&>1 || :;
endef

help:
	@echo "$(BOLD)Usage:$(NC)"
	@echo "  make [$(GREEN)target$(NC)]"
	@echo ""
	@echo "$(BOLD)Targets:$(NC)"
	@echo "  $(GREEN)all$(NC)            : Build all Docker images"
	@echo "  $(GREEN)base$(NC)           : Build the base Docker image ($(BASEIMAGE))"
	@echo "  $(GREEN)workspace$(NC)      : Build the workspace Docker image ($(WORKSPACEIMAGE))"
	@echo "  $(GREEN)golang-img$(NC)     : Build the Golang Docker image (golang-$(IMAGESUFFIX))"
	@echo "  $(GREEN)python-img$(NC)     : Build the Python Docker image (python-$(IMAGESUFFIX))"
	@echo "  $(GREEN)node-img$(NC)       : Build the Node.js Docker image (node-$(IMAGESUFFIX))"
	@echo "  $(GREEN)ts-img$(NC)         : Build the TypeScript Docker image (ts-$(IMAGESUFFIX))"
	@echo "  $(GREEN)flutter-img$(NC)    : Build the Flutter Docker image (flutter-$(IMAGESUFFIX))"
	@echo "  $(GREEN)clean$(NC)          : Remove all built Docker images"
	@echo "  $(GREEN)help$(NC)           : Show this help message"

base: dockerfiles/base $(SOURCES) clean-base
	$(call build_image,$(BASEIMAGE),dockerfiles/base)

clean-base:
	$(call clean_image,$(BASEIMAGE))

workspace: dockerfiles/workspace $(SOURCES) clean-workspace
	$(call build_image,$(WORKSPACEIMAGE),dockerfiles/workspace)

clean-workspace:
	$(call clean_image,$(WORKSPACEIMAGE))

golang-img: dockerfiles/golang $(SOURCES) clean-golang
	$(call build_image,golang-$(IMAGESUFFIX),dockerfiles/golang)

clean-golang:
	$(call clean_image,golang-$(IMAGESUFFIX))

python-img: dockerfiles/python $(SOURCES) clean-python
	$(call build_image,python-$(IMAGESUFFIX),dockerfiles/python)

clean-python:
	$(call clean_image,python-$(IMAGESUFFIX))

node-img: dockerfiles/node $(SOURCES) clean-node
	$(call build_image,node-$(IMAGESUFFIX),dockerfiles/node)

clean-node:
	$(call clean_image,node-$(IMAGESUFFIX))

ts-img: dockerfiles/ts $(SOURCES) clean-ts
	$(call build_image,ts-$(IMAGESUFFIX),dockerfiles/ts)

clean-ts:
	$(call clean_image,ts-$(IMAGESUFFIX))

flutter-img: dockerfiles/flutter $(SOURCES) clean-flutter
	$(call build_image,flutter-$(IMAGESUFFIX),dockerfiles/flutter)

clean-flutter:
	$(call clean_image,flutter-$(IMAGESUFFIX))

all: clean base workspace golang-img python-img node-img ts-img flutter-img

clean: clean-base clean-workspace clean-golang clean-python clean-node clean-ts clean-flutter

