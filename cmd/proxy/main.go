package main

import (
	"log"
	"net/http"

	"os"

	"github.com/LoSquadrato/crawltrip/internal/config"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

// TODO:
// - decouple request handling and database operations using NATS

func main() {

	const port = config.ProxyPort
	const targetHost = config.ProxyTargetHost

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v\n", err)
	}

	secretKey := os.Getenv("ADMIN_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("ADMIN_SECRET_KEY is not set in environment variables")
	}

	log.Printf("Starting proxy server on :%s\n", port)

	// Connect to NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	newProxy, err := NewProxy(targetHost)
	if err != nil {
		log.Fatalf("Error creating proxy: %v", err)
	}

	proxyConfig := &ProxyConfig{
		proxy:     newProxy,
		nc:        nc,
		secretKey: secretKey,
	}

	mux := http.NewServeMux()
	mux.Handle("/", proxyConfig.middlewareFilter(newProxy))
	mux.Handle("/proxy/ping", proxyConfig.middlewareAdminAuth(http.HandlerFunc(proxyConfig.pingDatabaseHandler)))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	s.ListenAndServe()
}
