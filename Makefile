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

db-server:
	@if command -v cockroach > /dev/null 2>&1; then \
		cockroach start-single-node --insecure \
			--store=/var/lib/cockroach \
			--listen-addr=localhost:26257; \
	else \
		@echo "missing CockroachDB binary."; \
	fi

.PHONY: run-cdb-migrations migrate-check-deps check-dsn-env

run-cdb-migrations: migrate-check-deps check-dsn-env
	migrate -source file://linkgraph/store/cdb/migrations -database  '$(subst postgresql,cockroach,${CDB_DSN}) up'


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

define dsn_missing_err
CDB_DSN envvar is undefined. To run the migrations this envvar
must point to a cockroach db instance. For example, if you are
running a local cockroachdb (with --insecure) and have created
a database called 'linkgraph' you can define the envvar by 
running:

export CDB_DSN='postgresql://root@localhost:26257/linkgraph?sslmode=disable'

endef
export dsn_missing_err

check-dsn-env:
ifndef CDB_DSN
	$(error ${dsn_missing_err})
endif