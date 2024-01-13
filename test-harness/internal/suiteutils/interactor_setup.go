package suiteutils

import (
	"context"

	. "github.com/onsi/gomega"
	"github.com/pseudo-su/golang-temporal-service-template/test-harness/internal/interactors"
)

func NewPubsubInteractor(ctx context.Context, config *interactors.PubsubInteractorInput) *interactors.PubsubInteractor {
	client, err := interactors.NewPubsubInteractor(ctx, config)
	Expect(err).ToNot(HaveOccurred())

	return client
}
