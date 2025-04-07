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
	@if ! command -v golangci-lint > /dev/null 2>&1; then \
		echo "[go get] installing golangci-lint"; \
		GO111MODULE=on go get -u github.com/golangci/golangci-lint/cmd/golangci-lint; \
	fi

db_server:
	@if command -v cockroach > /dev/null 2>&1; then \
		cockroach start-single-node --insecure \
			--store=/var/lib/cockroach \
			--listen-addr=localhost:26257; \
	else \
		@echo "missing CockroachDB binary."; \
	fi

.PHONY: migrate-check-deps

migrate-check-deps:
	@if [ -z `which migrate` ]; then \
		echo "[go get] installing golang-migrate cmd with cockroachdb support"; \
		if [ "${GO111MODULE}" = "off" ]; then \
			echo "[go get] installing github.com/golang-migrate/migrate/cmd/migrate"; \
			go get -tags 'cockroachdb postgres' -u github.com/golang-migrate/migrate/cmd/migrate; \
			go install -tags 'cockroach postgres' github.com/golang-migrate/migrate/cmd/migrate; \
		else \
			echo "[go get] installing github.com/golang-migrate/migrate/v4/cmd/migrate"; \
			go get -tags 'cockroachdb postgres' -u github.com/golang-migrate/migrate/v4/cmd/migrate; \
			go install -tags 'cockroach postgres' github.com/golang-migrate/migrate/v4/cmd/migrate; \
		fi \
	fi