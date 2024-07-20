package workflowutils

import (
	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
)

func GenerateUuid(ctx workflow.Context) uuid.UUID {
	encodedUUID := workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
		return uuid.New().String()
	})

	var uid string
	_ = encodedUUID.Get(&uid)
	return uuid.MustParse(uid)
}
