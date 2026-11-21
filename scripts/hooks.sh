#!/bin/sh
set -eu
unset CDPATH
cd "$(dirname "$0")/.."
if git config --get core.hooksPath >/dev/null 2>&1; then
    printf '%s\n' 'FIX: custom core.hooksPath is configured; have its owner integrate scripts/pre-push without replacing existing hooks.' >&2
    exit 1
elif [ "$?" -ne 1 ]; then
    printf '%s\n' 'FIX: cannot read Git hook configuration.' >&2
    exit 1
fi
common_dir=$(git rev-parse --git-common-dir)
common_dir=$(cd -P "$common_dir" && pwd -P)
reported_hook=$(git rev-parse --git-path hooks/pre-push)
reported_dir=$(dirname "$reported_hook")
if [ "${reported_dir##*/}" != hooks ] ||
    ! reported_parent=$(cd -P "$(dirname "$reported_dir")" && pwd -P) ||
    [ "$reported_parent" != "$common_dir" ] ||
    [ -L "$common_dir/hooks" ]; then
    printf '%s\n' 'FIX: Git hooks directory must be inside the common Git directory and must not be a symlink.' >&2
    exit 1
fi
hook_dir="$common_dir/hooks"
hook="$hook_dir/pre-push"
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
