### Test

TEST_UNIT_OPTS := -timeout 10s -count=1 -parallel=50
TEST_UNIT_TARGETS := $(APP_TARGETS)

LABEL_FILTER ?=

BLACKBOX_FILTER ?= ${LABEL_FILTER}
TEST_INTEGRATION_BLACKBOX_FILTER ?= ${BLACKBOX_FILTER}
TEST_INTEGRATION_BLACKBOX_OPTS := -p --succinct --timeout=4m --label-filter="${TEST_INTEGRATION_BLACKBOX_FILTER}"
TEST_INTEGRATION_BLACKBOX_REPORT_OPTS := -cover -coverprofile=test.integration.blackbox.coverage.out -coverpkg=./... -output-dir=./reports --junit-report=test.integration.blackbox.results.junit.xml --json-report=test.integration.blackbox.results.json
TEST_INTEGRATION_BLACKBOX_TARGETS := ./test-harness/suites/blackbox/...

ENV_FILTER ?= ${LABEL_FILTER}
TEST_ENV_LABEL_FILTER ?= ${ENV_FILTER}
TEST_ENV_OPTS := -p --succinct --timeout=10m --label-filter="${TEST_ENV_LABEL_FILTER}"
TEST_ENV_REPORT_OPTS := -cover -coverprofile=test.env.coverage.out -coverpkg=./... -output-dir=./reports --junit-report=test.env.results.junit.xml --json-report=test.env.results.json
TEST_ENV_TARGETS := ./test-harness/suites/env/...

test: test.unit test.integration
.PHONY: test

## Run unit tests
test.unit: test.unit.go
.PHONY: test.unit

## Run unit tests
test.unit.report: test.unit.go.report
.PHONY: test.unit.report

## Run Go unit tests
test.unit.go:
	go test ${TEST_UNIT_OPTS} ${TEST_UNIT_TARGETS}
.PHONY: test.unit.go

## Run unit tests and generate reports
test.unit.go.report:
	go test ${TEST_UNIT_OPTS} -coverprofile=reports/test.unit.coverage.out -coverpkg=./... ${TEST_UNIT_TARGETS} 2>&1 | ./tools/go-junit-report -set-exit-code -iocopy -out=./reports/test.unit.results.junit.xml
	go tool cover -html=reports/test.unit.coverage.out -o reports/test.unit.coverage.html
.PHONY: test.unit.go.report

## Run all integration tests
test.integration: test.integration.blackbox
.PHONY: test.integration

## Run all integration tests and generate reports
test.integration.report: test.integration.blackbox.report
.PHONY: test.integration.report

## Run blackbox integration tests
test.integration.blackbox:
	./tools/ginkgo ${TEST_INTEGRATION_BLACKBOX_OPTS} ${TEST_INTEGRATION_BLACKBOX_TARGETS}
.PHONY: test.integration.blackbox

## Run blackbox integration tests and generate reports
test.integration.blackbox.report:
	./tools/ginkgo ${TEST_INTEGRATION_BLACKBOX_OPTS} ${TEST_INTEGRATION_BLACKBOX_REPORT_OPTS} ${TEST_INTEGRATION_BLACKBOX_TARGETS}
	go tool cover -html=reports/test.integration.blackbox.coverage.out -o reports/test.integration.blackbox.coverage.html
.PHONY: test.integration.blackbox.report

## Run environment tests
test.env:
	if [ -z "${ENV}" ]; then echo "ENV environment variable not set"; exit 1; fi
	ENV=${ENV} ./tools/ginkgo ${TEST_ENV_OPTS} ${TEST_ENV_TARGETS}
.PHONY: test.env

## Run environment tests and generate reports
test.env.report:
	if [ -z "${ENV}" ]; then echo "ENV environment variable not set"; exit 1; fi
	ENV=${ENV} ./tools/ginkgo ${TEST_ENV_OPTS} ${TEST_ENV_REPORT_OPTS} ${TEST_ENV_TARGETS}
	go tool cover -html=reports/test.env.coverage.out -o reports/test.env.coverage.html
.PHONY: test.env.report
