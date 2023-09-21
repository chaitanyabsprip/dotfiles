#!/bin/sh

_depends() { have "$1" ||
	{ echo "${0##*/} depends on $1, please install and try again." && exit 1; }; }

prereqs() {
	runas="$1"
	case "$(uname)" in
	Darwin)
		_depends brew
		_depends xcode-select
		xcode-select --install || true
		brew install ninja gettext cmake
		;;
	Linux)
		match='s/ID="\{0,1\}\(.*\)"\{0,1\}/\1/'
		distro="$(grep ^ID= /etc/os-release | sed "$match")"
		case "$distro" in
		alpine) apk add build-base cmake coreutils curl unzip gettext-tiny-dev ;;
			# CentOS / RHEL / Fedora -- sometime later
		arch) $runas pacman -S base-devel cmake unzip ninja curl ;;
		opensuse*) $runas zypper install ninja cmake gcc-c++ gettext-tools curl ;;
		ubuntu) $runas apt-get -y --no-install-recommends install ninja-build gettext cmake unzip curl make ;;
		esac
		;;
	*) echo "${RED}Unknown or unsupported OS${NC}" && return ;;
	esac
}

install_neovim() {
	have nvim && exit
	runas="$([ ! "$(whoami)" = root ] && echo "sudo")"
	[ -n "$runas" ] && _depends sudo
	prereqs "$runas"
	mkdir ~/programs && {
		cd ~/programs || { echo "${RED}Failed to install neovim${NC}" && return; }
	}
	git clone https://github.com/neovim/neovim
	git checkout nightly
	cd neovim && make CMAKE_BUILD_TYPE=Release
	# shellcheck disable=2086
	have $runas || echo "Requires $runas to install"
	$runas make install
	rm -r ~/programs/neovim || true
	echo "${GREEN}neovim successfully installed${NC}"
}

install_neovim
