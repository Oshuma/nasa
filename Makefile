GOCACHE ?= $(CURDIR)/.cache

.PHONY: test

# test runs the full go test suite with verbose output, using a local cache directory
# to avoid sandbox permission issues.
test:
	@mkdir -p $(GOCACHE)
	@GOCACHE=$(GOCACHE) go test -v ./...
