#!/bin/sh
set -u

failed=0
for tool in git make sh awk cat chmod cmp dirname id ln mkdir mktemp rmdir rm uname; do
    if command -v "$tool" >/dev/null 2>&1; then
        printf 'PASS: %s available\n' "$tool"
    else
        printf 'FIX: provide %s\n' "$tool"
        failed=1
    fi
done

if sh scripts/check-toolchain.sh; then
    printf '%s\n' 'PASS: non-root Linux and pinned Go'
    root=$(GOTOOLCHAIN=local "${GO:-go}" env GOROOT 2>/dev/null)
    if [ -n "$root" ] && [ -x "$root/bin/gofmt" ]; then
        printf '%s\n' 'PASS: matching gofmt available'
    else
        printf '%s\n' 'FIX: repair the Go installation to provide matching gofmt'
        failed=1
    fi
    compiler=$(GOTOOLCHAIN=local "${GO:-go}" env CC 2>/dev/null)
    # Go allows arguments after the compiler executable.
    if [ -n "$compiler" ] && command -v "${compiler%% *}" >/dev/null 2>&1; then
        printf '%s\n' 'PASS: C compiler available for race tests'
    else
        printf '%s\n' 'FIX: provide the C compiler for the race tests'
        failed=1
    fi
else
    failed=1
fi

printf '%s\n' 'PASS: no required environment variables, mounted volumes, services or network access for the source gate'
printf '%s\n' 'INFO: GO and CC are optional toolchain overrides'
sh scripts/hooks.sh check || failed=1
exit "$failed"
