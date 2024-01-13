package suiteutils

import (
	"context"
)

func InitInteractAsTestHarness(ctx context.Context) {
	if TestsuiteConfig.EnvName() == "devstack" {
		InitDevstackInteractor(ctx)
	}
}
