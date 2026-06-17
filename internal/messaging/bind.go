package messaging

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

// DeclareAndBindStream creates or updates a JetStream stream with the given configuration
func DeclareAndBindStream(js jetstream.JetStream, config jetstream.StreamConfig) (jetstream.Stream, error) {
	// Create or update the stream for crawl messages
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	stream, err := js.CreateOrUpdateStream(
		ctx,
		config,
	)
	if err != nil {
		log.Printf("Failed to create or update stream: %v\n", err)
		return nil, err
	}
	return stream, nil
}
