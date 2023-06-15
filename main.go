package main

// kill %1 %2 %3 ; sleep 3 ; go run main.go TLD 8081 & go run main.go sub1 3000 & go run main.go sub2 3001 &
import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

//
// --------------------------------------------------------------------------------------
//                              Modify HTTP Header
// --------------------------------------------------------------------------------------

const addOriginHeader = false                    // add Access-Control header to HTTP response
var AllowOrigin string = "http://localhost:8081" // Choose a Access-Control origin header
//var AllowOrigin string = "http://localhost:3000"
//var AllowOrigin string = "http://localhost:3001"
//var AllowOrigin string = "http://localhost:222"
//var AllowOrigin string = "*"

type Message struct {
	Text string `json:"text"`
}

var Name string = ""
var Port string = "80"

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Need 2 args, Name and PortNumber")

	}
	Name = os.Args[1]
	Port = os.Args[2]
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", addHeaders(fs))
	http.HandleFunc("/get-json", jsonhandler)

	log.Print(Name + " Listening on :" + Port + "...")
	err := http.ListenAndServe(":"+Port, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func WriteACHeader(w http.ResponseWriter, AllowOrigin string) {
	if addOriginHeader {
		//w.Header().Add("X-Frame-Options", "GOFORIT")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Access-Control-Allow-Origin", AllowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
}

func addHeaders(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Add("X-Frame-Options", "GOFORIT")
		WriteACHeader(w, AllowOrigin)
		fs.ServeHTTP(w, r)
	}

}
func jsonhandler(w http.ResponseWriter, r *http.Request) {
	// Create a sample message
	message := Message{Text: "ThisPasswordIsSecretFor:" + Name}

	// Convert the message to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	WriteACHeader(w, AllowOrigin)

	// Write the JSON data to the response body
	w.Write(jsonData)
}
