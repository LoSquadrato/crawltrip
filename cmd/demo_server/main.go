package main

import (
	"log"
	"net/http"
)

// This is a simple backend server that we will use to test our proxy server.
// It listens on port 8081 and responds with "Hello, World!" to any request.

func main() {

	log.Println("Starting backend server on :8081")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s\n", r.Method, r.URL.Path)
		log.Printf("Proto: %s\n", r.Proto)
		log.Printf("Headers: %v\n", r.Header)
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		const page = `<html>
			<head></head>
			<body>
				<p> This is a backend server. </p>
			</body>
			</html>
			`
		w.Write([]byte(page))
	})

	log.Fatal(http.ListenAndServe(":8081", mux))
}
