package main

import (
	"net/http"
)

func main() {
	// Define a handler function
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			return
		}

		// Handle the actual request
		// Your code logic here...

		// Send a response
		w.Write([]byte("Hello, World!"))
	}

	// Create a new server and register the handler function
	http.HandleFunc("/", handler)

	// Start the server on port 8080
	http.ListenAndServe(":8080", nil)
}
