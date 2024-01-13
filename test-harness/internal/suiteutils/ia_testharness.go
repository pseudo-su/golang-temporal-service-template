package suiteutils

import (
	. "github.com/onsi/gomega"
	"github.com/pseudo-su/golang-temporal-service-template/test-harness/internal/interactors"
)

type InteractAsTestHarnessWorkflows struct{}

type InteractAsTestHarness struct {
	Devstack  *interactors.DevstackInteractor
	Workflows InteractAsTestHarnessWorkflows
}

func (interact InteractAsTestHarness) IsSetup() bool {
	interact.DevstackInteractorShouldBeSetup()

	return true
}

func (interact InteractAsTestHarness) DevstackInteractorShouldBeSetup() {
	if TestsuiteConfig.EnvName() == "devstack" {
		Expect(interact.Devstack).ToNot(BeNil())
		return
	}

	Expect(Iat.Devstack).To(BeNil())
}
