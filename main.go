package main

/*
                               CORS Lab

This program is a CORS security lab the allows users to explore both the client and server
side of CORS. Users can manipulate both the Client JavaScript and Server GO HTTP header CORS
attributes and view what the results are.

This is done thru the use of a main program and two iframes, all in different origins via
unique port numbers. The main program is a web server that forks to sub-processes all with
different port numbers which makes the document.location.origin unique. Users can then switch
between the 3 JS contexts (main, iframe1, iframe2) and view how CORS impacts accessing data
from the different origins. Users can do this with both

- getelementbyid to retireve data between origins
- postmessage between iframes to retrieve data

Usage:

    Start main process: go run main.go  HTML-Window-Name UniquePortNumber

Example:

Start 3 web servers with different port numbers. 1 main webserver and 2 others for the iframes

// Unix
kill $(jobs -p) ; sleep 3 ; go run main.go TLD 8081 & go run main.go iframe1 3000 & go run main.go iframe2 3001 &

// Windows
get-job| stop-job | remove-job ; go run main.go TLD 8081 & go run main.go iframes 3000 & go run main.go iframes 3001 &


once the main and 2 iframes are started, browse them. you have to pass the url that the frames will use to render their pages

http://localhost:8081/?iframeurl=iframes.html

press on the "Send Message" buttons to send messages postMessages between the
iframes and the parent, all in different origins.

Using Chrome Inspector or Firefox Debugger, to monitor the network traffic and
inspect the header fields for the CORS headers.const

You can then modify the

*/

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

const addOriginHeader = true // add Access-Control header to HTTP response
var AllowOrigin string = "*" // Choose a Access-Control origin header
//var AllowOrigin string = "http://localhost:8081"
//var AllowOrigin string = "http://localhost:3000"
//var AllowOrigin string = "http://localhost:3001"
//var AllowOrigin string = "http://localhost:222"

type Message struct {
	Text string `json:"text"`
}

var Name string = ""
var Port string = "80"

// --------------------------------------------------------------------------------------
//                              MAIN: is just a web server
// --------------------------------------------------------------------------------------
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

// --------------------------------------------------------------------------------------
//                              Add CORS headers to HTTP responses
// --------------------------------------------------------------------------------------

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

// --------------------------------------------------------------------------------------
//      Add message handler to recieve postMessages from iframes and return data
// --------------------------------------------------------------------------------------

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
