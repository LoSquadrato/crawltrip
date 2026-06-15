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

	"github.com/nats-io/nats.go"
)

type ProxyConfig struct {
	proxy     *httputil.ReverseProxy
	nc        *nats.Conn
	secretKey string
}

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(url), nil
}

// TODO:
// - add message queue to decouple request handling and database operations

func (p *ProxyConfig) middlewareFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")
		if err := Publisher(r, p.nc, userAgent); err != nil {
			log.Printf("Error publishing request: %v\n", err)
		}
		next.ServeHTTP(w, r)
	})
}

func (p *ProxyConfig) middlewareAdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretKey := r.Header.Get("X-Admin-Secret")
		if secretKey != p.secretKey {
			log.Printf("Unauthorized access attempt with secret key: %s\n", secretKey)
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Publisher(r *http.Request, nc *nats.Conn, userAgent string) error {
	if userAgent == "" {
		log.Printf("No User-Agent header found\n")
		return nil
	}
	if classificator.Classificator(userAgent) {
		log.Printf("Bot detected: %s\n", userAgent)
		// TODO: parse request return []byte
		rawreq, err := classificator.ParseRequest(r)
		if err != nil {
			log.Printf("Error parsing request: %v\n", err)
			return err
		}
		msg := &messaging.Message{
			Subject: config.CrawlSubjectPrefix + ".request",
			Data:    rawreq,
		}
		if err := messaging.PublishMessage(nc, msg); err != nil {
			log.Printf("Error publishing request: %v\n", err)
			return err
		}
		log.Printf("Request details: Method=%s, URL=%s, RemoteAddr=%s\n", rawreq.Method, rawreq.URL, rawreq.RemoteAddr)
	}
	return nil
}

func (p *ProxyConfig) pingDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Message string `json:"message"`
	}{
		Message: "ping",
	}
	msg := &messaging.Message{
		Subject: config.CrawlSubjectPrefix + ".ping",
		Data:    data,
	}
	if err := messaging.PublishMessage(p.nc, msg); err != nil {
		log.Printf("Error publishing ping message: %v\n", err)
	}
	log.Println("Ping request sent to database")
	utils.RespondWithJSON(w, http.StatusOK, data)
}
