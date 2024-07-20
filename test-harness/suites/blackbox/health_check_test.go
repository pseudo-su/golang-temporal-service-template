package blackbox_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/connection"
	. "github.com/pseudo-su/golang-temporal-service-template/test-harness/internal/suiteutils"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	grpc_deephealth_v1 "github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/deephealth/v1"
)

var _ = Describe("readiness", LabelDeliverable(DeliverableServiceAccess), func() {
	var healthCheckClient healthpb.HealthClient
	var conn *grpc.ClientConn
	var err error

	BeforeEach(func(ctx SpecContext) {
		conn, err = connection.GRPCDialContext(ctx, TestsuiteConfig.ServiceFrontdoorApiGrpcUri(), TestsuiteConfig.ServiceFrontdoorApiInsecure)
		Expect(err).ToNot(HaveOccurred())
	})

	It("readiness endpoint", func(ctx SpecContext) {
		healthCheckClient = healthpb.NewHealthClient(conn)
		response, err := healthCheckClient.Check(ctx, &healthpb.HealthCheckRequest{})

		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())
		Expect(response.Status).To(Equal(healthpb.HealthCheckResponse_SERVING))
	})

	It("deep health", func(ctx SpecContext) {
		client := grpc_deephealth_v1.NewDeepHealthClient(conn)
		response, err := client.Check(ctx, &grpc_deephealth_v1.DeepHealthCheckRequest{})
		Expect(err).ToNot(HaveOccurred())
		Expect(response).ToNot(BeNil())
		Expect(response.HealthState).To(Equal(grpc_deephealth_v1.DeepHealthCheckResponse_HEALTH_STATE_OK))
	})
})
