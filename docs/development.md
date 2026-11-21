# Development workflow

Run `make help` to find the entry points and `make help-all` for the full list. Use `make hooks` to install the local main-push guard. If you have a custom Git hook, ask its owner to integrate `scripts/pre-push` into the existing hook.

Run `make doctor` for read-only development prerequisite checks. Run `make verify` on a non-root Linux builder with the exact Go version in `go.mod` and a C compiler for race detection. The gate uses the standard library and runs offline. Maintainers route source gates through the approved builder workflow.

`make build` puts the candidate command-line tools under `build/bin`. Run `build/bin/zeo-release --help` on Linux for release-preflight usage. See [the installation requirements](../REQUIREMENTS.md) for the remaining work toward an integrated distribution. The development doctor covers toolchain and hook availability.
