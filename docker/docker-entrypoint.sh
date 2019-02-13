#!/usr/bin/env sh

set -e

# first arg is `-config` or `--some-option`
if [ "${1#-}" != "$1" ]; then
	set -- kratgo "$@"
fi

exec "$@"
