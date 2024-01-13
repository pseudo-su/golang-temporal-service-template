APP_PATHS :=  ./modules/service-pkg ./modules/service-cli ./modules/frontdoor ./modules/worker ./modules/mocks
APP_TARGETS := ./modules/service-pkg/... ./modules/service-cli/... ./modules/frontdoor/... ./modules/worker/... ./modules/mocks/...
DEV_PATHS :=  ./modules/testing-tools ./test-harness
DEV_TARGETS :=  ./modules/testing-tools/... ./test-harness/...
ALL_PATHS := $(APP_PATHS) $(DEV_PATHS)
ALL_TARGETS := $(APP_TARGETS) $(DEV_TARGETS)

include ./_help.mk
include ./tools/tools.mk
include ./_dev.mk
include ./_test.mk
include ./devstack/devstack.mk
