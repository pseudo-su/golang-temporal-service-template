tools/device-info.cfg: tools/device-info.sh
	./tools/device-info.sh > tools/device-info.cfg

tools/actionlint: tools/tools.cfg
	. ./tools/tools.cfg && env GOBIN=$${PWD}/tools go install github.com/rhysd/actionlint/cmd/actionlint@v$${actionlint}

tools/api-linter: tools/tools.cfg
	. ./tools/tools.cfg && env GOBIN=$${PWD}/tools go install github.com/googleapis/api-linter/cmd/api-linter@v$${api_linter_googleapis}

tools/buf: tools/tools.cfg
	. ./tools/tools.cfg && env GOBIN=$${PWD}/tools go install github.com/bufbuild/buf/cmd/buf@v$${buf}

tools/editorconfig-checker: tools/tools.cfg
	. ./tools/tools.cfg && env GOBIN=$${PWD}/tools go install github.com/editorconfig-checker/editorconfig-checker/v2/cmd/editorconfig-checker@latest
	# . ./tools/tools.cfg && env GOBIN=$${PWD}/tools GOPROXY=direct go install github.com/editorconfig-checker/editorconfig-checker/v2/cmd/editorconfig-checker@$${editorconfig_checker}

tools/ginkgo: tools/tools.cfg
	. ./tools/tools.cfg && env GOBIN=$${PWD}/tools go install github.com/onsi/ginkgo/v2/ginkgo@v$${ginkgo}

tools/go-junit-report: tools/tools.cfg
	. ./tools/tools.cfg && env GOBIN=$${PWD}/tools go install github.com/jstemmer/go-junit-report/v2@v$${go_junit_report}

tools/godotenv: tools/tools.cfg
	. ./tools/tools.cfg && env GOBIN=$${PWD}/tools go install github.com/joho/godotenv/cmd/godotenv@v$${godotenv}

tools/goimports: tools/tools.cfg
	. ./tools/tools.cfg && env GOBIN=$${PWD}/tools go install golang.org/x/tools/cmd/goimports@v$${goimports}

tools/golangci-lint: tools/tools.cfg
	. ./tools/tools.cfg && curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./tools v$${golangci_lint}

tools/grpcurl: tools/tools.cfg
	. ./tools/tools.cfg && env GOBIN=$${PWD}/tools go install github.com/fullstorydev/grpcurl/cmd/grpcurl@v$${grpcurl}

tools/mockery: tools/tools.cfg
	. ./tools/tools.cfg && env GOBIN=$${PWD}/tools go install github.com/vektra/mockery/v2@v$${mockery}

tools/plantuml.jar: tools/tools.cfg
	. ./tools/tools.cfg && curl -sfL https://github.com/plantuml/plantuml/releases/download/v$${plantuml}/plantuml.jar > ./tools/plantuml.jar

tools/protoc-gen-buf-breaking: tools/tools.cfg
	. ./tools/tools.cfg && env GOBIN=$${PWD}/tools go install github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking@v$${protoc_gen_buf_breaking}

tools/protoc-gen-buf-lint: tools/tools.cfg
	. ./tools/tools.cfg && env GOBIN=$${PWD}/tools go install github.com/bufbuild/buf/cmd/protoc-gen-buf-lint@v$${protoc_gen_buf_lint}

tools/protoc-gen-connect-go: tools/tools.cfg
	. ./tools/tools.cfg && env GOBIN=$${PWD}/tools go install connectrpc.com/connect/cmd/protoc-gen-connect-go@v$${protoc_gen_connect_go}

tools/protoc-gen-gapi-lint:
	. ./tools/tools.cfg && env GOBIN=$${PWD}/tools go install github.com/protoc-extensions/protoc-gen-gapi-lint/cmd/protoc-gen-gapi-lint@v$${protoc_gen_gapi_lint}

tools/protoc-gen-go-grpc: tools/tools.cfg
	. ./tools/tools.cfg && env GOBIN=$${PWD}/tools go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v$${protoc_gen_go_grpc}

tools/protoc-gen-go: tools/tools.cfg
	. ./tools/tools.cfg && env GOBIN=$${PWD}/tools go install google.golang.org/protobuf/cmd/protoc-gen-go@v$${protoc_gen_go}

tools/temporal: tools/tools.cfg tools/device-info.cfg
	. ./tools/tools.cfg && . ./tools/device-info.cfg && curl -sfL https://github.com/temporalio/cli/releases/download/v$${temporal_cli}/temporal_cli_$${temporal_cli}_$${device_platform}_$${device_architecture}.tar.gz > tools/temporal.tar.gz
	tar -xvzf ./tools/temporal.tar.gz temporal
	rm ./tools/temporal.tar.gz
	mv -f temporal ./tools
	touch tools/temporal

tools/workflowcheck: tools/tools.cfg
	. ./tools/tools.cfg && env GOBIN=$${PWD}/tools go install go.temporal.io/sdk/contrib/tools/workflowcheck@v$${workflowcheck}

tools/protoc: tools/tools.cfg tools/device-info.cfg
	. ./tools/tools.cfg && . ./tools/device-info.cfg && curl -sfL https://github.com/protocolbuffers/protobuf/releases/download/v$${protoc}/protoc-$${protoc}-$${protoc_platform}-$${protoc_architecture}.zip > tools/protoc.zip
	unzip -o ./tools/protoc.zip  bin/protoc -d ./tools
	rm ./tools/protoc.zip
	mv -f ./tools/bin/protoc ./tools/protoc
	rmdir ./tools/bin
	touch tools/protoc
