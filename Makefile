GO ?= go

.PHONY: verify format-check vet test test-bootstrap
verify: format-check vet test

format-check:
	@files="$$(gofmt -l cmd internal)" && test -z "$$files"

vet:
	$(GO) vet ./...

test:
	$(GO) test -race -count=1 ./...

test-bootstrap:
	$(GO) test -race -count=1 ./cmd/zeo-target-bootstrap
