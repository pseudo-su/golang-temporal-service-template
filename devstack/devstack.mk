### Devstack
DEVSTACK_TARGET ?= all
DEVSTACK_DOCKER_COMPOSE := docker-compose -f ./devstack/docker-compose.yaml
DEVSTACK_DOCKER_COMPOSE_INSTRUMENTED := docker-compose -f ./devstack/docker-compose.yaml  -f ./devstack/instrumented.docker-compose.yaml

## Start the devstack
devstack.start:
	$(DEVSTACK_DOCKER_COMPOSE) --verbose up --build -d --remove-orphans $(DEVSTACK_TARGET)
.PHONY: devstack.start

## Start the devstack with instrumented binaries for application components
devstack.start-instrumented:
	git clean -f -x "./devstack/components/**/coverage/*"
	$(DEVSTACK_DOCKER_COMPOSE_INSTRUMENTED) up --build -d --remove-orphans $(DEVSTACK_TARGET)
.PHONY: devstack.start-instrumented

# Restart all devstack containers
devstack.restart:
	$(DEVSTACK_DOCKER_COMPOSE) restart
.PHONY: devstack.restart

# Restart devstack containers for application component
devstack.restart.components:
	./devstack/docker-compose-restart-services.sh component_
.PHONY: devstack.restart.components

# Restart devstack containers for application dependencies
devstack.restart.deps:
	./devstack/docker-compose-restart-services.sh dep_
.PHONY: devstack.restart.deps

# Restart devstack containers for development tools
devstack.restart.tools:
	./devstack/docker-compose-restart-services.sh tool_
.PHONY: devstack.restart.tools

## Stop the devstack
devstack.stop:
	$(DEVSTACK_DOCKER_COMPOSE) down --remove-orphans
.PHONY: devstack.stop

## Clean/reset the devstack
devstack.clean:
	$(DEVSTACK_DOCKER_COMPOSE) down --remove-orphans --volumes --rmi local
.PHONY: devstack.clean

## Reload the devstack
devstack.reload: devstack.stop devstack.start
.PHONY: devstack.reload

## Clean/reset and restart the devstack
devstack.recreate: devstack.clean devstack.start
.PHONY: devstack.recreate

## show status
devstack.status:
	$(DEVSTACK_DOCKER_COMPOSE) ps
.PHONY: devstack.status

## Capture devstack coverage reports from application components
devstack.capture-coverage-reports: devstack.restart.components
	go tool covdata textfmt -i=devstack/components/frontdoor/coverage -o reports/test.integration.frontdoor.coverage.out
	go tool covdata textfmt -i=devstack/components/worker/coverage -o reports/test.integration.worker.coverage.out

	go tool cover -html=reports/test.integration.frontdoor.coverage.out -o reports/test.integration.frontdoor.coverage.html
	go tool cover -html=reports/test.integration.worker.coverage.out -o reports/test.integration.worker.coverage.html
.PHONY: devstack.capture-coverage-reports

## show logs
devstack.logs:
	$(DEVSTACK_DOCKER_COMPOSE) logs --follow
.PHONY: devstack.logs
