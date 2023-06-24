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
get-job| stop-job | remove-job ; go run main.go TLD 8081 & go run main.go iframes1 3000 & go run main.go iframes2 3001 &


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
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

//
// --------------------------------------------------------------------------------------
//                              Modify HTTP Header
// --------------------------------------------------------------------------------------

const addOriginHeader = true // add Access-Control header to HTTP response
//var AllowOrigin string = "*" // Choose a Access-Control origin header
//var AllowOrigin string = "http://localhost:8081"
//var AllowOrigin string = "http://localhost:3000"
//var AllowOrigin string = "http://localhost:3001"
var AllowOrigin string = "http://localhost:222"

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

	// Setup all the paths to handle HTTP requests
	fs := http.FileServer(http.Dir("./static"))
	// handle default web requests as front end server for html pages
	http.Handle("/", addHeaders(fs))

	// handle login forms from JS Fetch requests
	http.HandleFunc("/login", loginHandler)
	// handle login forms from classic Form Post Submit
	http.HandleFunc("/classic-form-submit", classicFormSubmit)

	// test program: return secrets to client, see if they read it
	http.HandleFunc("/get-json", jsonhandler)

	// HTTP server
	go func() {
		log.Print(Name + " Listening on HTTP port:" + Port + "...")
		err := http.ListenAndServe(":"+Port, nil)
		if err != nil {
			log.Fatal("Error starting HTTP server: ", err)
		}
	}()

	// HTTPS server
	go func() {
		HttpsPort, err := strconv.Atoi(Port)
		HttpsPort += 300
		log.Print(Name + " Listening on HTTPS port:" + strconv.Itoa(HttpsPort) + "...")
		certFile := "./certificate.crt"
		keyFile := "./privatekey.key"
		err = http.ListenAndServeTLS(":"+strconv.Itoa(HttpsPort), certFile, keyFile, nil)
		if err != nil {
			log.Fatal("Error starting HTTPS server: ", err)
		}
	}()
	select {}
}

// --------------------------------------------------------------------------------------
//      Login Logic to handle passwords
// --------------------------------------------------------------------------------------
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	Success bool   `json:"success"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Perform authentication logic here

	// Check if the username and password are valid
	success := loginReq.Username == "admin" && loginReq.Password == "password"
	message := ""
	token := ""
	if success {
		message = "Login successful!"
		token = "12345"
	} else {
		message = "Invalid username or password"
	}

	loginResp := LoginResponse{
		Message: message,
		Token:   token,
		Success: success,
	}

	w.Header().Set("Content-Type", "application/json")
	WriteACHeader(w, AllowOrigin)
	if err := json.NewEncoder(w).Encode(loginResp); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}

func classicFormSubmit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		WriteACHeader(w, AllowOrigin)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Must use a POST, preferabally HTTPS"))
	} else {
		WriteACHeader(w, AllowOrigin)
		r.ParseForm()
		// logic part of log in
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}

// --------------------------------------------------------------------------------------
//      Add HTTP Request Handler to recieve GET /get-json request to return data to client
//      so see if client can read it cross-origin
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
		WriteACHeader(w, AllowOrigin)
		fs.ServeHTTP(w, r)
	}

}
