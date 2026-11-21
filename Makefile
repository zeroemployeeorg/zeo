.DEFAULT_GOAL := help
GO ?= go
export GO

.PHONY: help help-all doctor hooks verify lint format-check vet test test-bootstrap build clean

help: ## Show the first-time development guide
	@printf '%s\n' 'ZEO distribution assembly tools' '' \
		'  make doctor    Check the local development prerequisites' \
		'  make hooks     Install the lightweight main-push guard' \
		'  make verify    Run the source gate on a Linux builder' \
		'  make build     Compile the candidate tools on a Linux builder' \
		'  make clean     Remove build output' '' \
		'Run make help-all for all targets. See docs/development.md for the workflow.'

help-all: ## List every documented target
	@awk 'BEGIN { FS = ":.*## " } /^[[:alnum:]_-]+:.*## / { printf "  %-18s %s\n", $$1, $$2 }' Makefile

doctor: ## Inspect tools and the Git hook without changing anything
	@sh scripts/doctor.sh

hooks: ## Install the main-push guard while preserving existing custom hooks
	@sh scripts/hooks.sh install

verify: lint vet test ## Check formatting, shell syntax, Go diagnostics and race-enabled tests

lint: format-check ## Check Go formatting and shell syntax
	@for file in scripts/*.sh scripts/pre-push; do sh -n "$$file" || exit 1; done

format-check: ## Report Go files that need gofmt
	@sh scripts/check-toolchain.sh
	@root="$$(GOTOOLCHAIN=local "$(GO)" env GOROOT)" && \
		files="$$("$$root/bin/gofmt" -l cmd internal)" && \
		{ test -z "$$files" || { printf '%s\n' "$$files"; exit 1; }; }

vet: ## Check Go compiler and vet diagnostics
	@sh scripts/check-toolchain.sh
	GOTOOLCHAIN=local "$(GO)" vet ./...

test: ## Run all Go tests with race detection and without cached results
	@sh scripts/check-toolchain.sh
	GOTOOLCHAIN=local CGO_ENABLED=1 "$(GO)" test -race -count=1 ./...

test-bootstrap: ## Run the focused target-bootstrap source tests
	@sh scripts/check-toolchain.sh
	GOTOOLCHAIN=local CGO_ENABLED=1 "$(GO)" test -race -count=1 ./cmd/zeo-target-bootstrap

build: ## Compile the candidate CLIs into build/bin
	@sh scripts/check-toolchain.sh
	@test ! -L build && test ! -L build/bin
	@mkdir -p build/bin
	GOTOOLCHAIN=local "$(GO)" build -trimpath -o build/bin/ ./cmd/...

clean: ## Remove only the local build output directory
	@test ! -L build
	@rm -rf -- build
