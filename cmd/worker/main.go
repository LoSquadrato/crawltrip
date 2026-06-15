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

type WorkerConfig struct {
	dbClient *mongo.Client
	broker   *nats.Conn
}

// TODO:
// - add aknowledgement and retry mechanism for failed messages with NATS JetStream
func main() {
	// Connect to MongoDB
	dbClient, err := database.Connect(config.DbUri)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}
	defer database.Close(dbClient, context.TODO())

	// Connect to NATS server
	url := config.NatsUrl
	if url == "" {
		url = nats.DefaultURL
	}
	// Connect to default NATS server (nats://127.0.0.1:4222)
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	// Empty the connection pool and close the connection when done
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

	// Create a context with a timeout for stream creation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	streamName := config.CrawlStreamName

	// Create a stream for the crawl messages
	stream, err := newJS.Stream(ctx, streamName)
	if err == nil {
		log.Fatalf("Failed to create stream: %v\n", err)
	}

	log.Println("Connected to NATS server!")

	// Queue group subscription
	// TODO:
	// - move subject and queue group to config
	// - add a logic to handle multiple subscriptions and handlers
	consCtx, err := messaging.CreateAndConsumeJSON(
		stream,
		ctx,
		config.CrawlConsumerName,
		streamName,
		config.CrawlSubjectPrefix+".*",
		workerConfig.HandlerMsg,
	)
	if err != nil {
		log.Fatal(err)
	}

	// get info for logging purposes
	streamInfo, _ := stream.Info(ctx)

	log.Printf("%s: %s on stream for listening for message on subject: %s\n", streamInfo.Config.Name, config.CrawlConsumerName, config.CrawlSubjectPrefix+".*")

	// Wait for a signal to gracefully exit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Exiting...")

	// Drain the subscription before exiting
	consCtx.Drain()

	// Drain the connection to ensure all messages are processed before exiting
	nc.Drain()

	// Delete the stream and consumer before exiting
	log.Println("# Delete stream")
	newJS.DeleteStream(ctx, streamName)
}
