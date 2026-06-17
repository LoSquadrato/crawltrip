// This is the main entry point for the proxy server, which will start the server and listen for incoming requests.
// It will also connect to the NATS server and create a new reverse proxy to the target host.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"

	"github.com/LoSquadrato/crawltrip/internal/config"
)

type ProxyConfig struct {
	nc        *nats.Conn
	secretKey string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v\n", err)
	}

	// Get the secret key from environment variables
	secretKey := os.Getenv("ADMIN_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("ADMIN_SECRET_KEY is not set in environment variables")
	}

	log.Printf("Starting proxy server on :%s\n", config.ProxyPort)

	// Connect to NATS server
	nc, err := nats.Connect(config.NatsUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Drain()

	// Create and start new reverse proxy
	newProxy, err := NewProxy(config.ProxyTargetHost)
	if err != nil {
		log.Fatalf("Error creating proxy: %v", err)
	}

	// Create proxy configuration
	proxyConfig := &ProxyConfig{
		nc:        nc,
		secretKey: secretKey,
	}

	mux := http.NewServeMux()
	mux.Handle("/", proxyConfig.MiddlewareFilter(newProxy))
	mux.HandleFunc("/proxy/ping", proxyConfig.MiddlewareAdminAuth(proxyConfig.PingDatabaseHandler))

	s := &http.Server{
		Addr:    ":" + config.ProxyPort,
		Handler: mux,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
