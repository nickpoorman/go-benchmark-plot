GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST?=$(GO_CMD) test
GO_GET=$(GO_CMD) get
GO_RUN=${GO_CMD} run

DIST_DIR=bin

GO_SOURCES := $(shell find . -path -prune -o -name '*.go' -not -name '*_test.go')
SOURCES_NO_VENDOR := $(shell find . -path ./vendor -prune -o -name "*.go" -not -name '*_test.go' -print)
CMD_TOOLS := $(shell find ./cmd -path -prune -o -name "*.go" -not -name '*_test.go')

.PHONY: default clean test build fmt bench clean-cache

#
# Our default target, clean up, do our install, test, and build locally.
#
default: clean build

# Clean up after our install and build processes. Should get us back to as
# clean as possible.
#
clean:
	@for d in ./bin/*; do \
		if [ -f $$d ] ; then rm $$d ; fi \
	done

clean-cache: clean
	go clean -cache -testcache -modcache

#
# Do what we need to do to run our tests.
#
test: clean $(GO_SOURCES)
	$(GO_TEST) -count=1 -v $(GO_TEST_ARGS) ./...

#
# Build/compile our application.
#
build:
	@for d in ./cmd/*; do \
		echo "Building ${DIST_DIR}/`basename $$d`"; \
		${GO_BUILD} -o ${DIST_DIR}/`basename $$d` $$d; \
	done

#
# Format the sources.
#
fmt: $(SOURCES_NO_VENDOR)
	goimports -w $^

#
# Run the benchmarks for the tools.
#
bench: $(GO_SOURCES)
	$(GO_TEST) $(GO_TEST_ARGS) -bench=. -run=- ./...

