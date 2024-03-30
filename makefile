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

help:
	@echo "$(BOLD)Usage:$(NC)"
	@echo "  make [$(GREEN)target$(NC)]"
	@echo ""
	@echo "$(BOLD)Targets:$(NC)"
	@echo "  $(GREEN)base$(NC)           : Build the base Docker image ($(BASEIMAGE))"
	@echo "  $(GREEN)workspace$(NC)      : Build the workspace Docker image ($(WORKSPACEIMAGE))"
	@echo "  $(GREEN)python-img$(NC)     : Build the Python Docker image (python-$(IMAGESUFFIX))"
	@echo "  $(GREEN)ts-img$(NC)         : Build the TypeScript Docker image (ts-$(IMAGESUFFIX))"
	@echo "  $(GREEN)flutter-img$(NC)    : Build the Flutter Docker image (flutter-$(IMAGESUFFIX))"
	@echo "  $(GREEN)clean$(NC)          : Remove all built Docker images"
	@echo "  $(GREEN)help$(NC)           : Show this help message"

all: base workspace
	@mkdir ~/projects

base: dockerfiles/base $(SOURCES) clean
	@docker build -t "$(BASEIMAGE)" -f dockerfiles/base

workspace: dockerfiles/workspace $(SOURCES) clean
	@docker build -t "$(WORKSPACEIMAGE)" -f dockerfiles/workspace

python-img: clean $(SOURCES)
	@docker build -t "python-$(IMAGESUFFIX)" -f dockerfiles/python

ts-img: clean $(SOURCES)
	@docker build -t "ts-$(IMAGESUFFIX)" -f dockerfiles/typescript

flutter-img: clean $(SOURCES)
	@docker build -t "flutter-$(IMAGESUFFIX)" -f dockerfiles/flutter

clean:
	@docker images -a | grep -q "$(BASEIMAGE)" && docker rmi "$(BASEIMAGE)"
	@docker images -a | grep -q "$(WORKSPACEIMAGE)" && docker rmi "$(WORKSPACEIMAGE)"
	@docker images -a | grep "-$(IMAGESUFFIX)$" | xargs docker rmi
