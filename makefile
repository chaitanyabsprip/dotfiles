SHELL = /bin/sh
CURRDIR := $(shell pwd)
export PATH := ${CURRDIR}/bin:$(PATH)

BOLD         := $(shell tput -Txterm bold)
RED          := $(shell tput -Txterm setaf 1)
GREEN        := $(shell tput -Txterm setaf 2)
YELLOW       := $(shell tput -Txterm setaf 3)
BLUE         := $(shell tput -Txterm setaf 6)
NC := $(shell tput -Txterm sgr0)

BINDIR = ${HOME}/.config/bin
MANPATH = /usr/local/share/man

all: zsh neovim tmux go

build: Dockerfile makefile setup
	@docker rmi dotfiles
	@docker build -t dotfiles .

ifneq ($(shell have zsh && echo 0 || echo 1),0)
zsh: ; @echo "${RED}○ Skipping zsh setup, not found${NC}"
else
zsh: --zap --starship --eza --fd --bat; @echo "${GREEN}✔ zsh setup completed${NC}"
endif

ifneq ($(shell have tmux && echo 0 || echo 1),0)
tmux: ; @echo "${RED}○ Skipping tmux setup, not found${NC}"
else
tmux: --gitmux --gitui --fzf; @echo "${GREEN}✔ tmux setup completed${NC}"
endif

ifeq ($(shell have nvim && echo 0 || echo 1),0)
neovim: ; @echo "${RED}○ Skipping nvim installation${NC}"
else
neovim:
	@bin/install-neovim.sh
	@echo "${GREEN}✔ nvim setup completed${NC}"
endif

go:
	@mkdir temp && { cd temp || :; }; \
		wget https://git.io/go-installer.sh; \
		chmod +x go-installer.sh; \
		bash go-installer.sh; \
		cd ..; \
		rm -r temp

ifeq ($(shell have gitmux && echo 0 || echo 1),0)
--gitmux: ; @echo "${BLUE}○ Skipping gitmux instalation, already present${NC}"
else
--gitmux:
	@echo "${YELLOW}◉ Installing gitmux${NC}"
	@spinner ghdl arl/gitmux
	@tar -xf gitmux-latest*
	@rm gitmux-latest* LICENSE README.md
	@mv gitmux ${BINDIR}/gitmux
	@echo "${GREEN}✔ gitmux installed${NC}"
endif

ifeq ($(shell have gitui && echo 0 || echo 1),0)
--gitui: ; @echo "${BLUE}○ Skipping gitui instalation, already present${NC}"
else
--gitui:
	@echo "${YELLOW}◉ Installing gitui${NC}"
	@mkdir temp && { cd temp || true; } && \
		spinner ghdl extrawurst/gitui && \
		tar -xf gitui-latest* && \
		mv gitui ~/.config/bin && \
		cd .. && rm -r temp;
	@echo "${GREEN}✔ gitui installed${NC}"
endif


ifeq ($(shell have fd && echo 0 || echo 1),0)
--fd: ; @echo "${BLUE}○ Skipping fd instalation, already present${NC}"
else
--fd:
	@echo "${YELLOW}◉ Installing fd${NC}"
ifeq ('Darwin', $(shell uname))
	@brew install bat
else
	@ghdl sharkdp/fd
	@install_path=${BINDIR}/fd && \
		tar -xf fd-latest.tar.gz && \
		asset=$$(find . -name 'fd-*' -type d) && \
		mv $$asset/fd $$install_path && \
		mv $$asset/fd.1 ${MANPATH}/ && \
		rm -r $$(find . -name 'fd-*')
	@echo "${GREEN}✔ fd installed${NC}"
endif
endif

ifeq ($(shell have bat && echo 0 || echo 1),0)
--bat: ; @echo "${BLUE}○ Skipping bat instalation, already present${NC}"
else
--bat:
	@echo "${YELLOW}◉ Installing bat${NC}"
ifeq ('Darwin', $(shell uname))
	@brew install bat
else
	@ghdl sharkdp/bat
	@install_path=${BINDIR}/bat && \
		tar -xf bat-latest.tar.gz && \
		asset=$$(find . -name 'bat-*' -type d) && \
		mv $$asset/bat $$install_path && \
		mv $$asset/bat.1 ${MANPATH}/ && \
		rm -r $$(find . -name 'bat-*')
	@echo "${GREEN}✔ bat installed${NC}"
endif
endif

ifeq ($(shell have yq && echo 0 || echo 1),0)
--yq: ; @echo "${BLUE}○ Skipping yq instalation, already present${NC}"
else
--yq:
	@echo "${YELLOW}◉ Installing yq${NC}"
	@ghdl mikefarah/yq
	@install_path="$$HOME"/.config/bin/yq && \
		mv yq-latest "$$install_path" && \
		chmod +x "$$install_path"
	@echo "${GREEN}✔ yq installed${NC}"
endif

ifeq ($(shell have fzf && echo 0 || echo 1),0)
--fzf: ; @echo "${BLUE}○ Skipping fzf instalation, already present${NC}"
else
--fzf:
	@printf "${YELLOW}◉ Installing fzf${NC} "
	@spinner curl -LsS https://raw.githubusercontent.com/junegunn/fzf/master/install -o install
	@spinner bash ./install --no-key-bindings --completion --no-update-rc
	@rm ./install
	@printf "\r\033[K${GREEN}✔ fzf installed${NC}\n"
endif

ifeq ($(shell have starship && echo 0 || echo 1),0)
--starship: ; @echo "${BLUE}○ Skipping starship instalation, already present${NC}"
else
--starship:
	@printf "${YELLOW}◉ Installing starship${NC} "
	@spinner curl -sS https://starship.rs/install.sh -o starship.sh
	@spinner sh starship.sh -y 1>/dev/null 2>&1
	@rm starship.sh;
	@printf "\r\033[K${GREEN}✔ starship installed${NC}\n"
endif

ifneq ($(wildcard ${HOME}/.local/share/zap),)
--zap: ; @echo "${BLUE}○ Skipping zap instalation, already present${NC}"
else
--zap:
	@echo "${YELLOW}◉ Installing zap${NC}"
	@curl -LsS https://raw.githubusercontent.com/zap-zsh/zap/master/install.zsh -o install
	@echo "N" | zsh install --branch release-v1 --keep 2>&1 /dev/null
	@rm install
	@echo "${GREEN}✔ zap installed${NC}"
endif

ifeq ($(shell have eza && echo 0 || echo 1),0)
--eza: ; @echo "${BLUE}○ Skipping eza instalation, already present${NC}"
else
--eza:
	@echo "${YELLOW}◉ Installing eza${NC}"
	@spinner ghdl eza-community/eza
	@tar -xf eza-latest*
	@mv eza ~/.config/bin
	@ln -sf ~/.config/bin/eza ~/.config/bin/exa
	@rm eza-latest*
	@echo "${GREEN}✔ eza installed${NC}"
endif
