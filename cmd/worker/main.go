// This is the main entry point for the worker, which will connect to the database and the NATS server,
// and start listening for messages on the configured subject.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LoSquadrato/crawltrip/internal/config"
	"github.com/LoSquadrato/crawltrip/internal/database"
	"github.com/LoSquadrato/crawltrip/internal/messaging"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// WorkerConfig holds the configuration for the worker, including the database client and any other dependencies
type WorkerConfig struct {
	dbClient *mongo.Client
	broker   *nats.Conn
}

func main() {
	// Connect to MongoDB
	dbClient, err := database.Connect(config.DbUri, 30*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}
	defer database.Close(dbClient, context.TODO())

	// Connect to NATS server
	nc, err := nats.Connect(config.NatsUrl)
	if err != nil {
		log.Fatalf("Failed to connect to NATS server: %v\n", err)
	}
	// Drain the connection to ensure all messages are processed before exiting
	defer nc.Drain()

	// Create worker configuration
	workerConfig := &WorkerConfig{
		dbClient: dbClient,
		broker:   nc,
	}

	// Create a JetStream context
	newJS, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("Failed to create JetStream context: %v\n", err)
	}

	// Create or update the stream for crawl messages
	stream, err := messaging.DeclareAndBindStream(
		newJS,
		jetstream.StreamConfig{
			Name:      config.CrawlStreamName,
			Subjects:  []string{config.CrawlSubjectPrefix + ".*"},
			Retention: jetstream.WorkQueuePolicy,
			Storage:   jetstream.MemoryStorage,
			MaxAge:    24 * time.Hour,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create or update stream: %v\n", err)
	}

	// Consume messages using the durable consumer
	consCtx, err := messaging.SubscribeAndProcess(
		stream,
		context.TODO(),
		config.CrawlConsumerName,
		config.CrawlSubjectPrefix+".*",
		workerConfig.HandlerMsgJSON,
	)
	if err != nil {
		log.Fatalf("Failed to create or update Worker stream: %v\n", err)
	}
	// Drain the subscription before exiting
	defer consCtx.Drain()

	log.Printf("%s: %s on stream for listening for message on subject: %s\n", config.CrawlStreamName, config.CrawlConsumerName, config.CrawlSubjectPrefix+".*")

	// Create or update the stream for DLQ messages
	_, err = messaging.DeclareAndBindStream(
		newJS,
		jetstream.StreamConfig{
			Name:      config.CrawlStreamName + "_DLQ",
			Subjects:  []string{"$JS.EVENT.ADVISORY.CONSUMER.MAX_DELIVERIES." + config.CrawlStreamName + ".*", "$JS.EVENT.ADVISORY.CONSUMER.MSG_TERMINATED." + config.CrawlStreamName + ".*"},
			Retention: jetstream.LimitsPolicy,
			Storage:   jetstream.FileStorage,
			MaxAge:    72 * time.Hour,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create or update DLQ stream: %v\n", err)
	}

	// Wait for a signal to gracefully exit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Exiting...")

	log.Println("Worker stopped.")
}
