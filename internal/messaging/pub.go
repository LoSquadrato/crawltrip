package messaging

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

// Message struct represents the structure of messages to be published to the queue
// can be extended with additional fields as needed, such as timestamp, priority, ContentType, etc.
type Message struct {
	Subject string
	Data    interface{}
}

// function to publish messages to the queue as JSON
func PublishMessage(nc *nats.Conn, msg *Message) error {
	msgData, err := json.Marshal(msg.Data)
	if err != nil {
		log.Printf("Failed to marshal task: %v", err)
		return err
	}
	if err := nc.Publish(msg.Subject, msgData); err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}
	return nil
}
