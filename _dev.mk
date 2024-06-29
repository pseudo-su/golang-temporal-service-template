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
	go get -u ./...
	go mod download
.PHONY: deps.app.update

## Tidy app dependencies
deps.app.tidy:
	go get ./...
	go mod download
.PHONY: deps.app.tidy

## Install tool dependencies
deps.tools.install: \
	tools/actionlint \
	tools/api-linter \
	tools/buf \
	tools/editorconfig-checker \
	tools/ginkgo \
	tools/go-junit-report \
	tools/godotenv \
	tools/goimports \
	tools/golangci-lint \
	tools/mockery \
	tools/plantuml.jar \
	tools/protoc \
	tools/protoc-gen-buf-breaking \
	tools/protoc-gen-buf-lint \
	tools/protoc-gen-connect-go \
	tools/protoc-gen-gapi-lint \
	tools/protoc-gen-go \
	tools/protoc-gen-go-grpc \
	tools/protoc-gen-grpc-gateway \
	tools/protoc-gen-openapi \
	tools/temporal \
	tools/workflowcheck
.PHONY: deps.tools.install

## Update tool dependencies
deps.tools.update: deps.tools.install
	@echo "WARNING: Any tool dependencies need to be updated manually"
.PHONY: deps.tools.update

## Tidy tool dependencies
deps.tools.tidy:
	@echo "WARNING: Any tool dependencies need to be tidied manually"
.PHONY: deps.tools.tidy

deps.tools.clean:
	git clean -f -x "./tools/*"
.PHONY: deps.tools.clean

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
ifeq ($(SKIP_VERIFY_BUF),true)
	echo "Skipping "verify.buf": SKIP_VERIFY_BUF=true";
else
	./tools/buf lint
endif
.PHONY: verify.buf

## Verify empty commit diff after codegen
verify.empty-git-diff:
	./scripts/verify-empty-git-diff.sh
.PHONY: verify.empty-git-diff

### Code generation

## Run all code generation
codegen: codegen.docs codegen.go codegen.mockery codegen.deps codegen.autoformat
.PHONY: codegen

## Run docs code generation
codegen.docs:
	./scripts/generate-docs.sh
.PHONY: codegen.docs

## Run Golang code generation
codegen.go:
	go generate ./...
.PHONY: codegen.go

## Run Golang code generation
codegen.mockery: tools/mockery
	./tools/mockery
.PHONY: codegen.mockery

codegen.autoformat:
	gofmt -s -w $(APP_PATHS) $(DEV_PATHS)
	./tools/goimports -w $(APP_PATHS) $(DEV_PATHS)
.PHONY: codegen.autoformat

codegen.deps: deps.tidy
.PHONY: codegen.deps
