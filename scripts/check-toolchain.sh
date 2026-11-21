#!/bin/sh
set -eu
if [ "$(uname -s)" != Linux ] || [ "$(id -u)" = 0 ]; then
    printf '%s\n' 'FIX: run this target as a non-root user on a Linux builder.' >&2
    exit 1
fi
required=$(awk '$1 == "go" { print "go" $2; exit }' go.mod)
actual=$(GOTOOLCHAIN=local "${GO:-go}" version 2>/dev/null | awk '{print $3}')
if [ -z "$required" ] || [ "$actual" != "$required" ]; then
    printf 'FIX: use the pinned Go toolchain (%s).\n' "$required" >&2
    exit 1
fi
