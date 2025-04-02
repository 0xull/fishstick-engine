.PHONY: deps test lint lint-check-deps

deps:
	@if [ "${GO111MODULE}" = "off"]; then \
		echo "[dep] fetching package dependencies";\
		go get -u github.com/golang/dep/cmd/dep;\
		dep ensure; \
	else \
		go get ./...; \
	fi

test:
	@echo "[go test] running tests and collecting coverage metrics"
	@go test -v -tags all_tests -race -coverprofile=coverag.txt -covermode=atomic ./...

lint: lint-check-deps
	@echo "[golangci-lint] linting sources";
	@golangci-lint run \
		-E misspell \
		-E golint \
		-E gofmt \
		-E unconvert \
		--exclude-use-default=false \
		-e SA1019 \
		--timeout 5m \
		./...

lint-check-deps:
	@if ! command -v golangci-lint > /dev/null 2>&1; then
		echo "[go get] installing golangci-lint"; \
		GO111MODULE=on go get -u github.com/golangci/golangci-lint/cmd/golangci-lint; \
	fi
