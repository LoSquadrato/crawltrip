package config

// This file contains all the configuration constants for the application.

// Broker key names
const (
	CrawlSubjectPrefix = "worker"
	WorkerQueueGroup   = "request_queue"
	CrawlStreamName    = "crawl_stream"
	CrawlConsumerName  = "crawl_consumer"
)

// Database configuration constants
const (
	DbUri        = "mongodb://localhost:27017"
	DbName       = "crawltrip"
	DbCollection = "requests"
)

// Proxy configuration constants
const (
	ProxyTargetHost = "http://localhost:8081"
	ProxyHost       = "localhost"
	ProxyPort       = "8080"
)

const (
	// NatsUrl is the URL of the NATS server.
	// if no url is provided, the default URL
	NatsUrl = ""
)

const (
	ApiUrl = "http://localhost:8082"
)

// NATS use the default URL so we don't need to define it here,
// but we can if we want to use a different URL.
// const NatsUrl = "nats://localhost:4222"
