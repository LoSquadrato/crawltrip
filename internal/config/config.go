package config

// This file contains all the configuration constants for the application.

// Broker key names
const (
	CrawlSubjectPrefix = "request_worker"
	CrawlStreamName    = "request_stream"
	CrawlConsumerName  = "request_consumer"
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
	NatsUrl = "http://localhost:4222"
)
