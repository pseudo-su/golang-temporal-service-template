package interactors

const (
	envVar_PUBSUB_EMULATOR_HOST = "PUBSUB_EMULATOR_HOST"
)

type DevstackInteractor struct {
	Pubsub *PubsubInteractor
}
