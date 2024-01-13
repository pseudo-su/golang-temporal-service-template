package suiteutils

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pseudo-su/golang-temporal-service-template/test-harness/internal/config"
	"github.com/pseudo-su/golang-temporal-service-template/test-harness/internal/interactors"
)

var TestsuiteConfig *config.TestsuiteEnvConfig

// Actors / Interactors
var Iap1 InteractAsPerson

// Interact as Test Harness
var Iat InteractAsTestHarness

func ResetIAP1() {
	Iap1 = InteractAsPerson{}
}

func SetDefaultENV(defaultEnvVal string) {
	env := os.Getenv("ENV")
	if env == "" {
		err := os.Setenv("ENV", defaultEnvVal)
		Expect(err).ToNot(HaveOccurred())
	}
}

func EnsureValidENV() {
	err := config.IsValidEnv(os.Getenv("ENV"))
	Expect(err).ToNot(HaveOccurred())
}

func InitConfigAndTestData() {
	var err error
	TestsuiteConfig, err = config.LoadEnvConfig()
	Expect(err).ToNot(HaveOccurred())
	TestsuiteConfig.Print(GinkgoWriter)
}

func InitDevstackInteractor(ctx context.Context) {
	Iat.Devstack = &interactors.DevstackInteractor{
		Pubsub: NewPubsubInteractor(ctx, &interactors.PubsubInteractorInput{
			TopicName:        "devstackcheck.topic",
			ProjectId:        "devstack",
			SubscriptionName: fmt.Sprintf("devstack.subscription.%s", strings.ReplaceAll(uuid.New().String(), "-", "")),
			EmulatorHost:     "localhost:8086",
		}),
	}

	listeningCtx := context.Background()

	err := Iat.Devstack.Pubsub.InitSubscription(listeningCtx)
	if err != nil {
		panic(fmt.Errorf("error initialising DevstackEvents interactor: %w", err))
	}
	Iat.Devstack.Pubsub.StartListening(listeningCtx)

	DeferCleanup(func() {
		if Iat.Devstack != nil {
			err := Iat.Devstack.Pubsub.DeleteSubscription(listeningCtx)
			if err != nil {
				panic(fmt.Errorf("error shutting down DevstackEvents subscription: %w", err))
			}
			listeningCtx.Done()
		}
	})
}
