package suiteutils

import (
	"fmt"
	"slices"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type ScenarioType string

func (t ScenarioType) LabelStr() string {
	return string(t)
}

var (
	ScenarioSmoke      = ScenarioType("smoke")
	ScenarioRegression = ScenarioType("regression")

	ScenarioTypes = []ScenarioType{
		ScenarioSmoke,
		ScenarioRegression,
	}
)

type Deliverable string

func (d Deliverable) LabelStr() string {
	return fmt.Sprintf("deliverable-%s", d)
}

type ActivationState string

func (ds ActivationState) LabelStr() string {
	return fmt.Sprintf("activation-%s", ds)
}

var (
	ActivationPending  = ActivationState("pending")
	ActivationComplete = ActivationState("complete")

	ActivationStates = []ActivationState{
		ActivationPending,
		ActivationComplete,
	}
)

var (
	DeliverableServiceAccess = Deliverable("ServiceAccess")
	DeliverableDay1          = Deliverable("Day1")

	DeliverablesPending = []Deliverable{
		DeliverableDay1,
	}
	DeliverablesActive = []Deliverable{
		DeliverableServiceAccess,
	}
)

// Labels
var (
	LabelSmoke      = Label(ScenarioSmoke.LabelStr())
	LabelRegression = Label(ScenarioRegression.LabelStr())
)

func LabelDeliverable(deliverables ...Deliverable) Labels {
	labelStrings := []string{}
	var countPending int
	var countActive int
	for _, d := range deliverables {
		// Label scenario with the deliverable
		labelStrings = append(labelStrings, d.LabelStr())
		// is deliverable pending or active
		var isPending = slices.Contains(DeliverablesPending, d)
		var isActive = slices.Contains(DeliverablesActive, d)
		if isPending {
			countPending++
		}
		if isActive {
			countActive++
		}
		if !isPending && !isActive {
			panic(fmt.Sprintf("Deliverable %s should be in either DeliverablesPending or DeliverablesActive", d))
		}
	}
	// By default assume/mark a test as only testing "active" deliverables
	var acitivationState = ActivationComplete
	if countPending > 0 {
		acitivationState = ActivationPending
	}
	labelStrings = append(labelStrings, acitivationState.LabelStr())
	return Label(labelStrings...)
}

func EnsureLabelledWithScenarioType(ctx SpecContext) {
	labelStrs := ctx.SpecReport().Labels()
	labels := Labels(labelStrs)
	matchAnyScenarioType := strings.Join(asStringSlice(ScenarioTypes...), " || ")
	Expect(
		labels.MatchesLabelFilter(matchAnyScenarioType),
	).To(
		BeTrue(),
		fmt.Sprintf(
			"missing scenario type: labels [%s] did not match filter (%s)",
			strings.Join(labelStrs, ","),
			matchAnyScenarioType,
		),
	)
}

func EnsureLabelledWithDeliverable(ctx SpecContext) {
	labelStrs := ctx.SpecReport().Labels()
	labels := Labels(labelStrs)
	allDeliverables := []Deliverable{}
	allDeliverables = append(allDeliverables, DeliverablesActive...)
	allDeliverables = append(allDeliverables, DeliverablesPending...)
	matchAnyDeliverable := strings.Join(asStringSlice(allDeliverables...), " || ")
	Expect(
		labels.MatchesLabelFilter(matchAnyDeliverable),
	).To(
		BeTrue(),
		fmt.Sprintf(
			"missing scenario type: labels [%s] did not match filter (%s)",
			strings.Join(labelStrs, ","),
			matchAnyDeliverable,
		),
	)
	matchAnyActivationState := strings.Join(asStringSlice(ActivationStates...), " || ")
	Expect(
		labels.MatchesLabelFilter(matchAnyActivationState),
	).To(
		BeTrue(),
		fmt.Sprintf(
			"missing scenario type: labels [%s] did not match filter (%s)",
			strings.Join(labelStrs, ","),
			matchAnyActivationState,
		),
	)
}

func asStringSlice[T interface{ LabelStr() string }](input ...T) []string {
	var result []string
	for _, item := range input {
		result = append(result, item.LabelStr())
	}
	return result
}
