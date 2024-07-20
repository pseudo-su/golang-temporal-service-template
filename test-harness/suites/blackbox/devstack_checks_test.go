package blackbox

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	grpc_deephealth_v1 "github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/deephealth/v1"
	"github.com/pseudo-su/golang-temporal-service-template/test-harness/internal/interactors"
	. "github.com/pseudo-su/golang-temporal-service-template/test-harness/internal/suiteutils"
)

var _ = Describe("devstack checks", LabelDeliverable(DeliverableServiceAccess), func() {
	var devstackEventInteractorWithFilters *interactors.PubsubEventConsumerInteractor

	BeforeEach(func(ctx SpecContext) {
		devstackEventInteractorWithFilters = Iat.Devstack.Pubsub.GetPubsubEventsInteractorWithFilters(
			interactors.HealthyResponses,
		)
	})

	AfterEach(func(ctx SpecContext) {
		devstackEventInteractorWithFilters.PrintEvents(GinkgoWriter)
	})

	It("pubsub should be able to publish messages", func(ctx SpecContext) {
		exists, err := Iat.Devstack.Pubsub.TopicExists(ctx)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(BeTrue())

		exists, err = Iat.Devstack.Pubsub.SubscriptionExists(ctx)
		Expect(err).ToNot(HaveOccurred())
		Expect(exists).To(BeTrue())

		err = Iat.Devstack.Pubsub.Publish(ctx, "test-event", &grpc_deephealth_v1.DeepHealthCheckResponse{
			HealthState: grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_OK, // ok
		})
		Expect(err).ToNot(HaveOccurred())

		err = Iat.Devstack.Pubsub.Publish(ctx, "test-event", &grpc_deephealth_v1.DeepHealthCheckResponse{
			HealthState: grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_FAIL, // ok
		})
		Expect(err).ToNot(HaveOccurred())

		err = devstackEventInteractorWithFilters.WaitForEvent(grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_OK)
		Expect(err).ToNot(HaveOccurred())
	})
})
