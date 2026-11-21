#!/bin/sh
set -eu
unset CDPATH GIT_DIR GIT_WORK_TREE GIT_COMMON_DIR GIT_CONFIG GIT_CONFIG_COUNT GIT_CONFIG_PARAMETERS
export GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null

scripts=$(cd "$(dirname "$0")" && pwd -P)
workspace=$(mktemp -d)
trap 'rm -rf -- "$workspace"' EXIT
trap 'exit 1' HUP INT TERM

new_repo() {
    git init -q --template= "$1"
    mkdir -p "$1/scripts"
    cp "$scripts/hooks.sh" "$scripts/pre-push" "$1/scripts/"
}

run_hook() (
    cd "$1"
    sh scripts/hooks.sh "$2"
)

expect_refusal() {
    expected=$1
    shift
    if output=$("$@" 2>&1); then
        printf 'FAIL: command unexpectedly succeeded (%s)\n' "$expected" >&2
        exit 1
    fi
    case "$output" in
        *"$expected"*) ;;
        *)
            printf 'FAIL: expected %s; got %s\n' "$expected" "$output" >&2
            exit 1
            ;;
    esac
}

normal="$workspace/normal"
linked="$workspace/linked"
new_repo "$normal"
git -C "$normal" add -- scripts/hooks.sh scripts/pre-push
git -C "$normal" -c user.name=Fixture -c user.email=fixture@example.invalid commit -qm fixture
git -C "$normal" worktree add --quiet --detach "$linked"
expect_refusal 'run make hooks' run_hook "$normal" check
run_hook "$linked" install
run_hook "$normal" check
run_hook "$linked" check

git -C "$normal" config --local core.hooksPath ''
expect_refusal 'custom core.hooksPath' run_hook "$normal" install
git -C "$normal" config --local --unset core.hooksPath

invalid="$workspace/invalid-config"
new_repo "$invalid"
printf '\n[invalid\n' >> "$invalid/.git/config"
expect_refusal 'cannot read Git hook configuration' run_hook "$invalid" install
test ! -e "$invalid/.git/hooks/pre-push"

symlink="$workspace/symlink"
external_hooks="$workspace/external/hooks"
new_repo "$symlink"
mkdir -p "$external_hooks"
ln -s "$external_hooks" "$symlink/.git/hooks"
expect_refusal 'Git hooks directory must be inside' run_hook "$symlink" check
expect_refusal 'Git hooks directory must be inside' run_hook "$symlink" install
test ! -e "$external_hooks/pre-push"

# Override the reported hook path while keeping the real Git common directory.
real_git=$(command -v git)
mkdir "$workspace/bin"
cat > "$workspace/bin/git" <<'GIT'
#!/bin/sh
if [ "$*" = 'rev-parse --git-path hooks/pre-push' ]; then
    printf '%s\n' "$FIXTURE_HOOK_PATH"
else
    exec "$FIXTURE_REAL_GIT" "$@"
fi
GIT
chmod +x "$workspace/bin/git"
(
    export PATH="$workspace/bin:$PATH"
    export FIXTURE_REAL_GIT="$real_git"
    export FIXTURE_HOOK_PATH="$external_hooks/pre-push"
    expect_refusal 'Git hooks directory must be inside' run_hook "$normal" install
)
test ! -e "$external_hooks/pre-push"
printf '%s\n' 'PASS: hook installation and Git directory refusal checks'
