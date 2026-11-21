#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
if git config --get core.hooksPath >/dev/null; then
    printf '%s\n' 'FIX: custom core.hooksPath is configured; have its owner integrate scripts/pre-push without replacing existing hooks.' >&2
    exit 1
fi
hook=$(git rev-parse --git-path hooks/pre-push)
case "${1:-}" in
check)
    if [ ! -L "$hook" ] && [ -f "$hook" ] && [ -x "$hook" ] && cmp -s scripts/pre-push "$hook"; then
        printf '%s\n' 'PASS: exact executable main-push guard installed'
    else
        printf '%s\n' 'FIX: run make hooks; an existing different hook requires owner integration'
        exit 1
    fi
    ;;
install)
    if [ -e "$hook" ] || [ -L "$hook" ]; then
        if [ -L "$hook" ] || [ ! -f "$hook" ] || ! cmp -s scripts/pre-push "$hook"; then
            printf '%s\n' 'FIX: existing hook differs; refusing to overwrite it. Have its owner integrate the main guard.' >&2
            exit 1
        fi
        chmod 755 "$hook"
    else
        hook_dir=$(dirname "$hook")
        mkdir -p "$hook_dir"
        staging=$(mktemp -d "$hook_dir/.pre-push.XXXXXXXX")
        trap 'rm -f "$staging/pre-push"; rmdir "$staging"' EXIT
        trap 'exit 1' HUP INT TERM
        cat scripts/pre-push > "$staging/pre-push"
        chmod 755 "$staging/pre-push"
        ln "$staging/pre-push" "$hook_dir"
    fi
    printf '%s\n' 'PASS: main-push guard installed; no build or network action'
    ;;
*) printf '%s\n' 'usage: sh scripts/hooks.sh check|install' >&2; exit 2 ;;
esac
