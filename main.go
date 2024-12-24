package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	// Define the proxy endpoint
	http.HandleFunc("/proxy", func(w http.ResponseWriter, r *http.Request) {
		// Get the "url" query parameter
		targetURL := r.URL.Query().Get("url")
		if targetURL == "" {
			http.Error(w, "URL is required as a query parameter", http.StatusBadRequest)
			return
		}

		// Validate the URL
		_, err := url.ParseRequestURI(targetURL)
		if err != nil {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		// Fetch the content from the target URL
		resp, err := http.Get(targetURL)
		if err != nil {
			http.Error(w, "Failed to fetch the URL", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Forward the response headers
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		// Forward the response status and body
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	// Start the server on port 8080
	log.Println("CORS Proxy server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
