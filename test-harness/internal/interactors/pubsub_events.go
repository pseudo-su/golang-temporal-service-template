package interactors

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"time"

	"cloud.google.com/go/pubsub"
	deephealth_v1 "github.com/pseudo-su/golang-temporal-service-template/modules/service-pkg/deephealth/v1"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/proto"
)

var (
	ErrSubscriptionDoesntExist = errors.New("subscription doesn't exist")
)

type PubsubEventConsumer interface {
	ReceiveEvent(event *PubsubEvent)
}

// TODO: change this to a more generic event/proto type
type PubsubEvent struct {
	Health *deephealth_v1.CheckResponse
}

type PubsubInteractor struct {
	client *pubsub.Client
	State  *pubsubInteractorState

	consumers      []PubsubEventConsumer
	EventsReceived []PubsubEvent
}

type pubsubInteractorState struct {
	topicName        string
	subscriptionName string

	Topic        *pubsub.Topic
	Subscription *pubsub.Subscription
}

type PubsubInteractorInput struct {
	ProjectId        string
	TopicName        string
	SubscriptionName string
	EmulatorHost     string

	Options []option.ClientOption
}

func NewPubsubInteractor(ctx context.Context, config *PubsubInteractorInput) (*PubsubInteractor, error) {
	origEmulatorHost := os.Getenv(envVar_PUBSUB_EMULATOR_HOST)
	if config.EmulatorHost != "" {
		if err := os.Setenv(envVar_PUBSUB_EMULATOR_HOST, config.EmulatorHost); err != nil {
			return nil, err
		}
	}

	pubsubClient, err := pubsub.NewClient(ctx, config.ProjectId, config.Options...)
	if err != nil {
		return nil, fmt.Errorf("error creating devstack event interactor: %w", err)
	}

	if err := os.Setenv(envVar_PUBSUB_EMULATOR_HOST, origEmulatorHost); err != nil {
		return nil, err
	}

	return &PubsubInteractor{
		client: pubsubClient,
		State: &pubsubInteractorState{
			topicName:        config.TopicName,
			subscriptionName: config.SubscriptionName,
		},
	}, nil
}

func (in *PubsubInteractor) StartListening(ctx context.Context) {
	go in.Receive(ctx)
}

func (in *PubsubInteractor) Publish(ctx context.Context, eventType string, message proto.Message) error {
	body, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %s", err.Error())
	}

	in.client.Topic(in.State.topicName).Publish(ctx, &pubsub.Message{
		Data: body,
	})

	return nil
}

func (in *PubsubInteractor) Receive(ctx context.Context) {
	err := in.State.Subscription.Receive(ctx, func(ctx context.Context, message *pubsub.Message) {
		msg := deephealth_v1.CheckResponse{}

		err := proto.Unmarshal(message.Data, &msg)
		if err != nil {
			fmt.Println("failed to unmarshal event from pubsub")
		}

		event := PubsubEvent{
			Health: &msg,
		}

		// send all received events to all consumers
		for _, consumer := range in.consumers {
			// apply filters at consumer level in specified ReceiveEvent fn
			consumer.ReceiveEvent(&event)
		}
	})
	if err != nil {
		fmt.Printf("failed to receive events from subscription '%s': %s", in.State.subscriptionName, err.Error())
	}
}

func (in *PubsubInteractor) TopicExists(ctx context.Context) (bool, error) {
	exists, err := in.client.Topic(in.State.topicName).Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking existence of topic: %w", err)
	}

	if !exists {
		return exists, fmt.Errorf("topic doesn't exist")
	}

	return exists, nil
}

func (in *PubsubInteractor) SubscriptionExists(ctx context.Context) (bool, error) {
	exists, err := in.client.Subscription(in.State.subscriptionName).Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking existence of subscription: %w", err)
	}

	if !exists {
		return exists, ErrSubscriptionDoesntExist
	}

	return exists, nil
}

func (in *PubsubInteractor) InitSubscription(ctx context.Context) error {
	topic := in.client.Topic(in.State.topicName)

	in.State.Topic = topic

	subscription, err := in.client.CreateSubscription(ctx, in.State.subscriptionName, pubsub.SubscriptionConfig{Topic: topic})
	if err != nil {
		return fmt.Errorf("cannot create devstack pubsub subscription: %w", err)
	}

	in.State.Subscription = subscription
	return nil
}

func (in *PubsubInteractor) DeleteSubscription(ctx context.Context) error {
	err := in.State.Subscription.Delete(ctx)
	if err != nil {
		return fmt.Errorf("cannot delete subscription: %w", err)
	}

	return nil
}

type PubsubEventFilter func(*PubsubEvent) error

type PubsubEventConsumerInteractor struct {
	eventsReceived   []*PubsubEvent
	eventFilters     []PubsubEventFilter
	timeout          time.Duration
	waitedEventIndex int
}

func (in *PubsubEventConsumerInteractor) ReceiveEvent(event *PubsubEvent) {
	for _, filter := range in.eventFilters {
		if err := filter(event); err != nil {
			continue
		}
		in.eventsReceived = append(in.eventsReceived, event)
	}
}

func (in *PubsubEventConsumerInteractor) PrintEvents(w io.Writer) {
	if len(in.eventsReceived) < 1 {
		fmt.Fprintln(w, "no events")
	}

	for _, event := range in.eventsReceived {
		fmt.Fprintf(w, "EVENT: %s\n", event.Health.Message)
	}

	fmt.Fprintf(w, "total: %d", len(in.eventsReceived))
}

func (in *PubsubInteractor) GetPubsubEventsInteractorWithFilters(filterFns ...PubsubEventFilter) *PubsubEventConsumerInteractor {
	consumer := &PubsubEventConsumerInteractor{
		eventsReceived: make([]*PubsubEvent, 0),
		eventFilters:   filterFns,
		timeout:        time.Second * 3,
	}

	for _, e := range in.EventsReceived {
		consumer.ReceiveEvent(&e)
	}

	in.consumers = append(in.consumers, consumer)
	return consumer
}

func HealthyResponses(event *PubsubEvent) error {
	if event.Health.Message != "1" {
		return fmt.Errorf("unhealthy state found: %s", event.Health.Message)
	}
	return nil
}

func (in *PubsubEventConsumerInteractor) confirmEventReceived(expectedHealthState string) error {
	for _, event := range in.eventsReceived[in.waitedEventIndex:] {
		in.waitedEventIndex++
		switch event.Health.Message {
		case expectedHealthState:
			if reflect.TypeOf(event.Health.Message) == reflect.TypeOf(expectedHealthState) {
				return nil
			}
		case "2":
			continue
		case "0":
			continue
		default:
			continue
		}
	}

	return fmt.Errorf("did not receive event: %s", expectedHealthState)
}

func (in *PubsubEventConsumerInteractor) WaitForEvent(expected string) error {
	giveUp := time.Now().Add(in.timeout)
	for {
		if in.confirmEventReceived(expected) == nil {
			return nil
		} else if time.Now().After(giveUp) {
			return fmt.Errorf("failed to receive specified event %+v in time", expected)
		}

		time.Sleep(10 * time.Millisecond)
	}
}
