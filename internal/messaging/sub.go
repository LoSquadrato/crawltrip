package messaging

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// AckAction represents the acknowledgment action to be taken after processing a message.
type AckAction int

const (
	Ack AckAction = iota
	NackDiscard
	NackRetry
	NackWithDelay
)

// ConsumeAndProcess subscribes to a NATS subject and processes incoming messages using the provided handler function.
// working with NATS Core Pub/Sub model, don't work with JetStream
func ConsumeAndProcess(nc *nats.Conn, subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
	sub, err := nc.Subscribe(subject, handler)
	if err != nil {
		log.Printf("Error subscribing to subject %s: %v\n", subject, err)
		return nil, err
	}
	log.Printf("Subscribed to subject: %s\n", subject)
	return sub, nil
}

// QueueSubscribeAndProcess subscribes to a NATS subject with a queue group and processes incoming messages using the provided handler function.
// working with NATS Core Pub/Sub model, don't work with JetStream
func QueueSubscribeAndProcess(nc *nats.Conn, subject, queue string, handler nats.MsgHandler) (*nats.Subscription, error) {
	sub, err := nc.QueueSubscribe(subject, queue, handler)
	if err != nil {
		log.Printf("Error subscribing to subject %s with queue %s: %v\n", subject, queue, err)
		return nil, err
	}
	log.Printf("Subscribed to subject: %s with queue: %s\n", subject, queue)
	return sub, nil
}

// Create a durable consumer for the crawl messages which body is a JSON and process them using the provided handler function.
func SubscribeAndProcess(
	js jetstream.Stream,
	ctx context.Context,
	consumerName,
	subject string,
	handler func(m jetstream.Msg) AckAction,
) (jetstream.ConsumeContext, error) {

	// create a durable consumer for the crawl messages
	cons, err := js.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:          consumerName,
		FilterSubject: subject,
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxDeliver:    5,
		BackOff:       []time.Duration{5 * time.Second, 10 * time.Second, 30 * time.Second},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create or update consumer: %w", err)
	}
	log.Println("Created consumer", cons.CachedInfo().Name)

	// consume messages using the durable consumer
	consumeContext, err := cons.Consume(func(m jetstream.Msg) {
		action := handler(m)
		switch action {
		case Ack:
			m.Ack()
		case NackDiscard:
			m.Term()
		case NackRetry:
			m.Nak()
		case NackWithDelay:
			m.NakWithDelay(10 * time.Second) // example delay, adjust as needed
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %w", err)
	}

	return consumeContext, nil
}
