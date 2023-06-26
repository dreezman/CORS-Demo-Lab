package main

/*
                               CORS Lab

This program is a CORS security lab the allows users to explore both the client and server
side of CORS. Users can manipulate both the Client JavaScript and Server GO HTTP header CORS
attributes and view what the results are.

See README.md for details

*/

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var Name string = ""
var Port string = "80"

// --------------------------------------------------------------------------------------
//                              MAIN: is just a web server handling requests
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

	// handle toggle CORS headers
	http.HandleFunc("/cors-toggle", corsToggle)

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
//                              Add CORS headers to HTTP responses
// --------------------------------------------------------------------------------------
var addOriginHeader = true   // add Access-Control header to HTTP response
var AllowOrigin string = "*" // Choose a Access-Control origin header
func WriteACHeader(w http.ResponseWriter, AllowOrigin string) {
	if addOriginHeader {
		//w.Header().Add("X-Frame-Options", "GOFORIT")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Access-Control-Allow-Origin", AllowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Expose-Headers", "*")
	}
}

func addHeaders(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteACHeader(w, AllowOrigin)
		fs.ServeHTTP(w, r)
	}

}

// Let browser set the CORS header setting
func corsToggle(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	addOriginHeader = true
	param1 := r.URL.Query().Get("corssettings")
	switch param1 {
	case "TurnCorsOff":
		addOriginHeader = false
	case "TurnCorsWildOn":
		AllowOrigin = "*"
	case "TurnCorsRandomOrigOn":
		AllowOrigin = "https://xyz.com:123"
	case "TurnCorsParenOrigOn":
		AllowOrigin = "http://localhost:8081"
	default:
		AllowOrigin = "*"
		addOriginHeader = true
	}

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

	if r.Method == http.MethodOptions {
		WriteACHeader(w, AllowOrigin)
		w.WriteHeader(http.StatusOK)
		return
	}

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
	// Add access token, append port number of the origin
	u, _ := url.Parse(r.Header["Origin"][0])
	_, port, _ := net.SplitHostPort(u.Host)
	var cookie http.Cookie
	if success {
		message = "Login successful!"
		token = "12345"
		if r.TLS == nil {
			cookieorigin := "AccessToken_" + string(port)
			cookie = http.Cookie{Name: cookieorigin, Value: "12345", Domain: "localhost", Secure: false, SameSite: http.SameSiteLaxMode}
		} else {
			intport, _ := strconv.Atoi(port) // tls port on cookie is 300 higher always
			intport += 300
			cookieorigin := "AccessToken_" + strconv.Itoa(intport)
			cookie = http.Cookie{Name: cookieorigin, Value: "12345", Domain: "localhost", Secure: true, SameSite: http.SameSiteNoneMode}
		}
		http.SetCookie(w, &cookie)
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
type Message struct {
	Text string `json:"text"`
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
