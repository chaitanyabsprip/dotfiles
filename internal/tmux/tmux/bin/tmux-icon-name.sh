#!/bin/sh

_depends() { type "$1" >/dev/null 2>&1 ||
	{ echo "${0##*/} depends on $1, please install and try again." &&
		exit 1; }; }

_depends yq

NAME="$1"
CURRENT_DIR=$(dirname "$(readlink -f "$0")")
DEFAULT_CONFIG="$CURRENT_DIR/icon_names.yml"

get_config_value() {
	key=$1
	value="$(yq "$key" "$DEFAULT_CONFIG")"
	echo "$value"
}

ICON="$(get_config_value ".icons.$NAME")"

if [ "$ICON" = "null" ]; then
	FALLBACK_ICON="$(get_config_value '.config.fallback-icon')"
	ICON="$FALLBACK_ICON"
fi

SHOW_NAME="$(get_config_value '.config.show-name')"
if [ "$SHOW_NAME" = true ]; then
	ICON="$ICON $NAME"
fi

echo "$ICON"
