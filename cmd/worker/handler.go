package main

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go/jetstream"

	"context"

	"github.com/LoSquadrato/crawltrip/internal/classificator"
	"github.com/LoSquadrato/crawltrip/internal/config"
	"github.com/LoSquadrato/crawltrip/internal/database"
	"github.com/LoSquadrato/crawltrip/internal/messaging"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	selectorRequest = "request"
	selectorPing    = "ping"
)

func (wc *WorkerConfig) HandlerMsg(m jetstream.Msg) messaging.AckAction {
	switch m.Subject() {
	case config.CrawlSubjectPrefix + "." + selectorRequest:
		return handlerSaveRequest(wc.dbClient, m.Data())
	case config.CrawlSubjectPrefix + "." + selectorPing:
		return handlerPingDatabase(wc.dbClient, m.Data())
	default:
		log.Printf("Received message on unknown subject: %s\n", m.Subject())
		return messaging.NackDiscard
	}
}

func handlerSaveRequest(client *mongo.Client, data []byte) messaging.AckAction {
	var rw classificator.RawRequest
	if err := json.Unmarshal(data, &rw); err != nil {
		log.Printf("Failed to unmarshal request: %v\n", err)
		return messaging.NackDiscard
	}
	if err := database.SaveRequest(context.TODO(), client, &rw); err != nil {
		log.Printf("Failed to save request: %v\n", err)
		return messaging.NackWithDelay
	} else {
		log.Printf("Request saved: %s\n", rw.URL)
	}
	return messaging.Ack
}

func handlerPingDatabase(client *mongo.Client, data []byte) messaging.AckAction {
	var payload struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("Failed to unmarshal ping message: %v\n", err)
		return messaging.NackDiscard
	}
	if err := database.Ping(client, context.TODO()); err != nil {
		log.Printf("Failed to ping database: %v\n", err)
		return messaging.NackWithDelay
	} else {
		log.Println("Database ping successful")
		return messaging.Ack
	}
}
