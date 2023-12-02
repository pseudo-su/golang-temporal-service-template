include tools/tools.mk

APP_PATHS := ./modules/frontdoor ./modules/worker ./modules/service-pkg
APP_TARGETS := ./modules/frontdoor/... ./modules/worker/... ./modules/service-pkg/...
DEV_PATHS := ./test-harness
DEV_TARGETS := ./test-harness/...
ALL_PATHS := $(APP_PATHS) $(DEV_PATHS)
ALL_TARGETS := $(APP_TARGETS) $(DEV_TARGETS)

### Dependencies - manage project and tool dependencies

## Install dependencies
deps.install: deps.tools.install deps.app.install
.PHONY: deps.install

## Update dependencies
deps.update: deps.tools.update deps.app.update
.PHONY: deps.update

## Tidy dependencies
deps.tidy: deps.tools.tidy deps.app.tidy
.PHONY: deps.tidy

## Install app dependencies
deps.app.install:
	# Golang dependencies (go.work, go.mod, go.sum)
	go mod download
.PHONY: deps.app.install

## Update app dependencies
deps.app.update:
	@for dir in $(ALL_PATHS); do \
					echo "Running go get in $$dir"; \
					(cd $$dir && go get -u ./... ); \
	done
	go work sync
.PHONY: deps.app.update

## Tidy app dependencies
deps.app.tidy:
	go work sync
.PHONY: deps.app.tidy

## Install tool dependencies
deps.tools.install: \
	tools/golangci-lint \
	tools/plantuml.jar \
	tools/godotenv \
	tools/grpcurl \
	tools/go-junit-report \
	tools/mockery \
	tools/goimports \
	tools/workflowcheck \
	tools/actionlint \
	tools/ginkgo \
	tools/editorconfig-checker \
	tools/buf \
	tools/protoc \
	tools/protoc-gen-go \
	tools/protoc-gen-go-grpc \
	tools/protoc-gen-grpc-gateway \
	tools/protoc-gen-buf-breaking \
	tools/protoc-gen-buf-lint \
	tools/protoc-gen-openapi \
	tools/protoc-gen-gapi-lint \
	tools/api-linter
.PHONY: deps.tools.install

## Update tool dependencies
deps.tools.update: deps.tools.install
	echo "WARNING: Any tool dependencies need to be updated manually"
.PHONY: deps.tools.update

## Tidy tool dependencies
deps.tools.tidy:
	echo "WARNING: Any tool dependencies need to be tidied manually"
.PHONY: deps.tools.tidy

deps.tools.clean:
	git clean -f -x "./tools/*"
.PHONY: deps.tools.clean

### Test
TEST_UNIT_OPTS := -timeout 10s -count=1 -parallel=50
TEST_UNIT_TARGETS := $(APP_TARGETS)

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

## Run blackbox integration tests (require devstack and app to be running)
test.integration.blackbox:
	echo "TODO: Add scenario-driven blackbox integration tests"
.PHONY: test.integration.blackbox

## Run smoke tests
test.env.smoke:
	if [ -z "${ENV}" ]; then echo "ENV environment variable not set"; exit 1; fi
	echo "TODO: Add env smoke tests"
.PHONY: test.env.smoke

### Verify - Code verifiation and Static analysis

## Run code verification
verify: verify.go verify.editorconfig verify.github-actions verify.temporal-workflows verify.buf
.PHONY: verify

## Run code verifiation and automatically apply fixes where possible
verify.fix: verify.go.fix verify.editorconfig verify.github-actions verify.temporal-workflows verify.buf
.PHONY: verify.fix

# Run static analysis on Golang code
verify.go:
	./tools/golangci-lint run -v $(ALL_TARGETS)
.PHONY: verify.go

# Run static analysis on Golang code and autofix where possible
verify.go.fix:
	./tools/golangci-lint run --fix $(ALL_TARGETS)
.PHONY: verify.go.fix

# Run static analysis on .editorconfig rules
verify.editorconfig:
	./tools/editorconfig-checker --exclude modules/service-pkg/openapi.yaml
.PHONY: verify.editorconfig

## Verify using temporal workflow static analysis
verify.github-actions:
	./tools/actionlint -shellcheck=
.PHONY: verify.github-actionlint

## Verify using temporal workflow static analysis
verify.temporal-workflows:
	./tools/workflowcheck -config workflowcheck.config.yaml $(APP_TARGETS)
.PHONY: verify.temporal-workflows

verify.buf:
	./tools/buf lint
.PHONY: verify.buf

## Verify empty commit diff after codegen
verify.empty-git-diff:
	./scripts/verify-empty-git-diff.sh
.PHONY: verify.empty-git-diff

### Code generation

## Run all code generation
codegen: codegen.docs codegen.go codegen.mockery codegen.autoformat
.PHONY: codegen

## Run docs code generation
codegen.docs:
	./scripts/generate-docs.sh
.PHONY: codegen.docs

## Run Golang code generation
codegen.go:
# Open issue for using go.work w/ `go generate ./...`
# https://github.com/golang/go/issues/56098
	@for dir in $(ALL_PATHS); do \
					echo "Running go generate in $$dir"; \
					(cd $$dir && go generate ./... ); \
	done
.PHONY: codegen.go

## Run Golang code generation
codegen.mockery: tools/mockery
	./tools/mockery
.PHONY: codegen.mockery

codegen.autoformat:
	go work sync

	gofmt -s -w $(APP_PATHS) $(DEV_PATHS)
	./tools/goimports -w $(APP_PATHS) $(DEV_PATHS)
.PHONY: codegen.autoformat

### Docker

## Remove dangling docker images
docker.remove-untagged:
	docker rmi $$(docker images -f dangling=true -q) 2>/dev/null || true
.PHONY: docker.remove-untagged
