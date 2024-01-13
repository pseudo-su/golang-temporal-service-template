package blackbox_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/pseudo-su/golang-temporal-service-template/test-harness/internal/suiteutils"
)

var _ = BeforeSuite(func(ctx context.Context) {
	SetDefaultENV("devstack")
	EnsureValidENV()
	InitConfigAndTestData()

	InitInteractAsTestHarness(ctx)
})

var _ = JustBeforeEach(func(ctx SpecContext) {
	// EnsureLabelledWithScenarioType(ctx)
	EnsureLabelledWithDeliverable(ctx)
})
