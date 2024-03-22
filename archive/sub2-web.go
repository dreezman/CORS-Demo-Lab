package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", addHeaders(fs))

	log.Print("SubWeb2 Listening on :3001...")
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal(err)
	}
}
func addHeaders(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Add("X-Frame-Options", "ALLOW")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		fs.ServeHTTP(w, r)
	}
}
