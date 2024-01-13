package env_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/pseudo-su/golang-temporal-service-template/test-harness/internal/suiteutils"
)

var _ = Describe("smoke health check", LabelSmoke, func() {
	It("example active feature", LabelDeliverable(DeliverableServiceAccess), func(ctx SpecContext) {
		var err error
		Expect(err).ToNot(HaveOccurred())
	})
	It("example pending deliverable", LabelDeliverable(DeliverableDay1), func(ctx SpecContext) {
		var err error
		Expect(err).ToNot(HaveOccurred())
	})
})
