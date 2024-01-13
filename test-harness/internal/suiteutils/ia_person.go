package suiteutils

type InteractAsPerson struct{}

func (interact InteractAsPerson) IsSetup() bool {
	return true
}
