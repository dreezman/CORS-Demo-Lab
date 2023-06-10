package main

// go run static-web.go TLD 8081 & go run static-web.go sub1 3000 & go run static-web.go sub2 3001 &
import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Need 2 args, Name and PortNumber")

	}
	Name := os.Args[1]
	Port := os.Args[2]
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", addHeaders(fs))
	http.HandleFunc("/get-json", jsonhandler)

	log.Print(Name + " Listening on :" + Port + "...")
	err := http.ListenAndServe(":"+Port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
func addHeaders(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Add("X-Frame-Options", "GOFORIT")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:222")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		fs.ServeHTTP(w, r)
	}

}
func jsonhandler(w http.ResponseWriter, r *http.Request) {
	// Create a sample message
	message := Message{Text: "Hello, World!"}

	// Convert the message to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response body
	w.Write(jsonData)
}
