// This file contains the main handler for the proxy server, which includes the middleware to filter out bot requests
// and the admin endpoint to ping the database.
package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/LoSquadrato/crawltrip/internal/classificator"
	"github.com/LoSquadrato/crawltrip/internal/config"
	"github.com/LoSquadrato/crawltrip/internal/messaging"
	"github.com/LoSquadrato/crawltrip/internal/utils"
)

// Create a new reverse proxy to the target host and return it as an http.Handler
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(url), nil
}

// MiddlewareFilter is a middleware that checks if the incoming request is from a bot or not, using the classificator package.
func (p *ProxyConfig) MiddlewareFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")
		if userAgent == "" {
			log.Printf("No User-Agent header found\n")
		}
		if classificator.IsCrawler(userAgent) {
			log.Printf("Bot detected: %s\n", userAgent)
			rawreq, err := classificator.ParseRequest(r)
			if err != nil {
				log.Printf("Error parsing request: %v\n", err)
			}
			err = messaging.PublishMessage(
				p.nc,
				&messaging.Message{
					Subject: config.CrawlSubjectPrefix + ".request",
					Data:    rawreq,
				})
			if err != nil {
				log.Printf("Error publishing request: %v\n", err)
			}
			log.Println("Request published to NATS")
			log.Printf("Request details: Method=%s, URL=%s, RemoteAddr=%s\n", rawreq.Method, rawreq.URL, rawreq.RemoteAddr)
		}
		next.ServeHTTP(w, r)
	})
}

// MiddlewareAdminAuth is a middleware that checks for the presence of a secret key in the request header
func (p *ProxyConfig) MiddlewareAdminAuth(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		secretKey := r.Header.Get("X-Admin-Secret")
		if secretKey != p.secretKey {
			log.Printf("Unauthorized access attempt with secret key: %s\n", secretKey)
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		next(w, r)
	}
}

// PingDatabaseHandler is an admin endpoint that sends a ping message to the database
// through NATS include every part of the messaging chain.
// Check every single instance (proxy, worker, NATS, database) if they are alive and responsive.
func (p *ProxyConfig) PingDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	err := messaging.PublishMessage(
		p.nc,
		&messaging.Message{
			Subject: config.CrawlSubjectPrefix + ".ping",
			Data:    "ping",
		})
	if err != nil {
		log.Printf("Error publishing ping message: %v\n", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to ping database")
		return
	}
	log.Println("Ping request sent to database")
	utils.RespondWithJSON(w, http.StatusOK, "Database ping request sent successfully")
}
