package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go/jetstream"

	"context"

	"github.com/LoSquadrato/crawltrip/internal/classificator"
	"github.com/LoSquadrato/crawltrip/internal/config"
	"github.com/LoSquadrato/crawltrip/internal/database"
	"github.com/LoSquadrato/crawltrip/internal/messaging"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// enum for the different selectors for the messages, to route them to the correct handler
const (
	selectorRequest = "request"
	selectorPing    = "ping"
)

// general handler for the worker, it will route the message to the correct handler based on the subject of the message
func (wc *WorkerConfig) HandlerMsgJSON(m jetstream.Msg) messaging.AckAction {
	switch m.Subject() {
	case config.CrawlSubjectPrefix + "." + selectorRequest:
		return handlerSaveRequestJSON(wc.dbClient, m.Data())
	case config.CrawlSubjectPrefix + "." + selectorPing:
		return handlerPingDatabaseJSON(wc.dbClient, m.Data())
	default:
		log.Printf("Received message on unknown subject: %s\n", m.Subject())
		return messaging.NackDiscard
	}
}

// handler for saving the request to the database
func handlerSaveRequestJSON(client *mongo.Client, data []byte) messaging.AckAction {
	var rw classificator.RawRequest
	if err := json.Unmarshal(data, &rw); err != nil {
		log.Printf("Failed to unmarshal request: %v\n", err)
		return messaging.NackDiscard
	}
	if err := database.SaveRequest(context.TODO(), client, rw); err != nil {
		log.Printf("Failed to save request: %v\n", err)
		return messaging.NackWithDelay
	} else {
		log.Printf("Request saved: %s\n", rw.URL)
	}
	return messaging.Ack
}

// handler for pinging the database to check if it's alive and responsive
func handlerPingDatabaseJSON(client *mongo.Client, data []byte) messaging.AckAction {
	payload := string(data)
	log.Printf("Received ping message with payload: %s\n", payload)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := database.Ping(client, ctx); err != nil {
		log.Printf("Failed to ping database: %v\n", err)
		return messaging.NackWithDelay
	} else {
		log.Println("Database ping successful")
		return messaging.Ack
	}
}
