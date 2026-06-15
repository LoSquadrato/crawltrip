package classificator

// This package provides a simple classificator to determine if a user agent is a bot or not.
// We start using only the crawler-user-agents library, but we can easily extend it to use other sources
// of information in the future, such as a database or an API.

import (
	"bytes"
	"io"
	"net/http"
	"time"

	//"github.com/google/uuid"
	agents "github.com/monperrus/crawler-user-agents"
)

func Classificator(userAgent string) bool {
	return agents.IsCrawler(userAgent)
}

type RawRequest struct {
	TimeStamp  time.Time           `json:"timestamp"`
	Method     string              `json:"method"`
	URL        string              `json:"url"`
	Proto      string              `json:"proto,omitempty"`
	Headers    map[string][]string `json:"headers,omitempty"`
	Body       []byte              `json:"body,omitempty"`
	RemoteAddr string              `json:"remote_addr,omitempty"`
	Host       string              `json:"host,omitempty"`
}

func ParseRequest(r *http.Request) (*RawRequest, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return &RawRequest{
		TimeStamp:  time.Now(),
		Method:     r.Method,
		URL:        r.URL.String(),
		Proto:      r.Proto,
		Headers:    map[string][]string(r.Header),
		Body:       bodyBytes,
		RemoteAddr: r.RemoteAddr,
		Host:       r.Host,
	}, nil
}
