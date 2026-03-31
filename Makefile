.PHONY: build run test test-integration lint vet fmt clean help serve list check

# Connection
DATABASE_URL ?= postgres://pgextras:pgextras@db:5432/pgextras_test?sslmode=disable

# Build
BINARY      := pgextras
BUILD_DIR   := bin
CMD_PKG     := ./cmd/pgextras

## help: Show this help message
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## //' | column -t -s ':'

## build: Compile the CLI binary to bin/
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) $(CMD_PKG)

## run: Build and run a query (usage: make run Q=cache_hit)
run: build
	DATABASE_URL=$(DATABASE_URL) $(BUILD_DIR)/$(BINARY) query $(Q)

## serve: Start the web dashboard on :8080
serve: build
	DATABASE_URL=$(DATABASE_URL) $(BUILD_DIR)/$(BINARY) serve --public --addr :8080

## list: List all available queries
list: build
	$(BUILD_DIR)/$(BINARY) list

## test: Run unit tests
test:
	go test -v -race -count=1 ./...

## test-integration: Run integration tests against the devcontainer database
test-integration:
	DATABASE_URL=$(DATABASE_URL) go test -v -race -count=1 -tags=integration ./...

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## vet: Run go vet
vet:
	go vet ./...

## fmt: Format all Go source files
fmt:
	gofmt -s -w .

## check: Run fmt, vet, lint, and tests
check: fmt vet lint test

## clean: Remove build artifacts
clean:
	rm -rf $(BUILD_DIR)

## diagnose: Run health checks against the database
diagnose: build
	DATABASE_URL=$(DATABASE_URL) $(BUILD_DIR)/$(BINARY) diagnose

## add-extensions: Install pg_stat_statements, pg_buffercache, sslinfo
add-extensions: build
	DATABASE_URL=$(DATABASE_URL) $(BUILD_DIR)/$(BINARY) query add_extensions
