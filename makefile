.PHONY: clean

BOLD         := $(shell tput -Txterm bold)
RED          := $(shell tput -Txterm setaf 1)
GREEN        := $(shell tput -Txterm setaf 2)
YELLOW       := $(shell tput -Txterm setaf 3)
BLUE         := $(shell tput -Txterm setaf 6)
NC           := $(shell tput -Txterm sgr0)

BASEIMAGE = clinic
WORKSPACEIMAGE = hospital
IMAGESUFFIX = ward
SOURCES := $(filter-out dockerfiles/%, $(wildcard */*))
SOURCES := $(filter-out $(wildcard $(shell find . -type l -printf '%P\n')), $(SOURCES))


help:
	@echo "$(BOLD)Usage:$(NC)"
	@echo "  make [$(GREEN)target$(NC)]"
	@echo ""
	@echo "$(BOLD)Targets:$(NC)"
	@echo "  $(GREEN)base$(NC)           : Build the base Docker image ($(BASEIMAGE))"
	@echo "  $(GREEN)workspace$(NC)      : Build the workspace Docker image ($(WORKSPACEIMAGE))"
	@echo "  $(GREEN)golang-img$(NC)     : Build the Golang Docker image (golang-$(IMAGESUFFIX))"
	@echo "  $(GREEN)python-img$(NC)     : Build the Python Docker image (python-$(IMAGESUFFIX))"
	@echo "  $(GREEN)node-img$(NC)       : Build the Node.js Docker image (node-$(IMAGESUFFIX))"
	@echo "  $(GREEN)ts-img$(NC)         : Build the TypeScript Docker image (ts-$(IMAGESUFFIX))"
	@echo "  $(GREEN)flutter-img$(NC)    : Build the Flutter Docker image (flutter-$(IMAGESUFFIX))"
	@echo "  $(GREEN)clean$(NC)          : Remove all built Docker images"
	@echo "  $(GREEN)help$(NC)           : Show this help message"

all: base workspace
	@mkdir ~/projects

base: dockerfiles/base $(SOURCES) clean-base
	@docker build $(if $(DEBUG),--progress=plain) -t "$(BASEIMAGE)" -t "chaitanyabsprip/$(BASEIMAGE)" -f dockerfiles/base .

clean-base:
	@docker images -a --format '{{.Repository}}' | grep -w "$(BASEIMAGE)" | \
		xargs -I {} docker rmi {} 2>/dev/null || :

workspace: dockerfiles/workspace $(SOURCES) clean-workspace
	@docker build$(if $(DEBUG),--progress=plain) -t "$(WORKSPACEIMAGE)" -t "chaitanyabsprip/$(WORKSPACEIMAGE)" -f dockerfiles/workspac .

clean-workspace:
	@docker images -a --format '{{.Repository}}' | grep -w  "$(WORKSPACEIMAGE)" | \
		xargs -I {} docker rmi {} 2>/dev/null || :

golang-img: dockerfiles/golang $(SOURCES) clean-golang
	@docker build $(if $(DEBUG),--progress=plain) -t "golang-$(IMAGESUFFIX)" -t "chaitanyabsprip/golang-$(IMAGESUFFIX)" -f dockerfiles/golang .

clean-golang:
	@docker images -a --format '{{.Repository}}' | grep -- "golang-$(IMAGESUFFIX)" | \
		xargs -I {} docker rmi {} 2>/dev/null || :

python-img: dockerfiles/python $(SOURCES) clean-python
	@docker build $(if $(DEBUG),--progress=plain) -t "python-$(IMAGESUFFIX)" -t "chaitanyabsprip/python-$(IMAGESUFFIX)" -f dockerfiles/python .

clean-python:
	@docker images -a --format '{{.Repository}}' | grep -- "python-$(IMAGESUFFIX)" | \
		xargs -I {} docker rmi {} 2>/dev/null || :

node-img: dockerfiles/node $(SOURCES) clean-node
	@docker build $(if $(DEBUG),--progress=plain) -t "node-$(IMAGESUFFIX)" -t "chaitanyabsprip/node-$(IMAGESUFFIX)" -f dockerfiles/node .

clean-node:
	@docker images -a --format '{{.Repository}}' | grep -- "node-$(IMAGESUFFIX)" | \
		xargs -I {} docker rmi {} 2>/dev/null || :

ts-img: dockerfiles/ts $(SOURCES) clean-ts
	@docker build $(if $(DEBUG),--progress=plain) -t "ts-$(IMAGESUFFIX)" -t "chaitanyabsprip/ts-$(IMAGESUFFIX)" -f dockerfiles/ts .

clean-ts:
	@docker images -a --format '{{.Repository}}' | grep -- "ts-$(IMAGESUFFIX)" | \
		xargs -I {} docker rmi {} 2>/dev/null || :

flutter-img: dockerfiles/flutter $(SOURCES) clean-flutter
	@docker build $(if $(DEBUG),--progress=plain) -t "flutter-$(IMAGESUFFIX)" -t "chaitanyabsprip/flutter-$(IMAGESUFFIX)" -f dockerfiles/flutter .

clean-flutter:
	@docker images -a --format '{{.Repository}}' | grep -- "flutter-$(IMAGESUFFIX)" | \
		xargs -I {} docker rmi {} 2>/dev/null || :

all: base workspace golang-img python-img node-img ts-img flutter-img

clean: clean-base clean-workspace clean-golang clean-python clean-node clean-ts clean-flutter
