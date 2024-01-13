### Devstack

DEVSTACK_TARGET ?= all
DEVSTACK_COMPOSE_TOOL ?= podman compose

DEVSTACK_COMPOSE := $(DEVSTACK_COMPOSE_TOOL) -f ./devstack/compose.yaml
DEVSTACK_COMPOSE_INSTRUMENTED := $(DEVSTACK_COMPOSE_TOOL) -f ./devstack/compose.yaml  -f ./devstack/compose.instrumented.yaml

## Pull all devstack containers
devstack.pull:
	$(DEVSTACK_COMPOSE) pull
.PHONY: devstack.pull

## Build all devstack containers
devstack.build:
	$(DEVSTACK_COMPOSE) build
.PHONY: devstack.build

## Start the devstack
devstack.start:
	$(DEVSTACK_COMPOSE) --verbose up --build -d --remove-orphans $(DEVSTACK_TARGET)
.PHONY: devstack.start

## Start the devstack with instrumented binaries for application components
devstack.start-instrumented:
	git clean -f -x "./devstack/components/**/coverage/*"
	$(DEVSTACK_COMPOSE_INSTRUMENTED) up --build -d --remove-orphans $(DEVSTACK_TARGET)
.PHONY: devstack.start-instrumented

# Restart all devstack containers
devstack.restart:
	$(DEVSTACK_COMPOSE) restart
.PHONY: devstack.restart

# Restart devstack containers for application component
devstack.restart.components:
	./devstack/compose-restart-services.sh component_
.PHONY: devstack.restart.components

# Restart devstack containers for application dependencies
devstack.restart.deps:
	./devstack/compose-restart-services.sh dep_
.PHONY: devstack.restart.deps

# Restart devstack containers for development tools
devstack.restart.tools:
	./devstack/compose-restart-services.sh tool_
.PHONY: devstack.restart.tools

## Stop the devstack
devstack.stop:
	$(DEVSTACK_COMPOSE) down --remove-orphans
.PHONY: devstack.stop

## Clean/reset the devstack
devstack.clean:
	$(DEVSTACK_COMPOSE) down --remove-orphans --volumes --rmi local
.PHONY: devstack.clean

## Reload the devstack
devstack.reload: devstack.stop devstack.start
.PHONY: devstack.reload

## Clean/reset and restart the devstack
devstack.recreate: devstack.clean devstack.start
.PHONY: devstack.recreate

## show status
devstack.status:
	$(DEVSTACK_COMPOSE) ps
.PHONY: devstack.status

## Capture devstack coverage reports from application components
devstack.capture-coverage-reports: devstack.restart.components
	if [ -z "${TEST_SUITE}" ]; then echo "TEST_SUITE environment variable not set"; exit 1; fi

	go tool covdata textfmt -i=devstack/components/frontdoor/coverage -o reports/${TEST_SUITE}.devstack.frontdoor.coverage.out
	go tool covdata textfmt -i=devstack/components/worker/coverage -o reports/${TEST_SUITE}.devstack.worker.coverage.out

	go tool cover -html=reports/${TEST_SUITE}.devstack.frontdoor.coverage.out -o reports/${TEST_SUITE}.devstack.frontdoor.coverage.html
	go tool cover -html=reports/${TEST_SUITE}.devstack.worker.coverage.out -o reports/${TEST_SUITE}.devstack.worker.coverage.html
.PHONY: devstack.capture-coverage-reports

### Devstack logs

LOG_TARGET ?= devstack-component_worker
DEVSTACK_LOG_TARGET ?= $(LOG_TARGET)

devstack.logs: devstack.logs.follow
.PHONY: devstack.logs

## Export devstack logs
devstack.logs.export:
	$(DEVSTACK_COMPOSE) logs
.PHONY: devstack.logs.export

## Add devstack logs to reports/
devstack.logs.report:
	$(DEVSTACK_COMPOSE) logs --no-color > reports/devstack.log
.PHONY: devstack.logs.report

# Follow devstack logs
devstack.logs.follow:
	$(DEVSTACK_COMPOSE) logs --follow
.PHONY: devstack.logs.follow

## Follow devstack filtered logs
devstack.logs.target:
	podman logs --follow $(shell podman ps -aqf "name=$(DEVSTACK_LOG_TARGET)")
.PHONY: devstack.logs.target
