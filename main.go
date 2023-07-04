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
	"strings"
)

var WebServerName string = ""
var WebServerHTTPPort string = "80"




// --------------------------------------------------------------------------------------
//
//	MAIN: is just a web server handling requests
//
// --------------------------------------------------------------------------------------
func main() {
	if len(os.Args) < 2 {
		log.Fatal("Need 2 args, WebServerName and PortNumber")

	}
	WebServerName = os.Args[1]
	WebServerHTTPPort = os.Args[2]

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

	// set cookies in response
	http.HandleFunc("/get-cookies", cookiehandler)

	// HTTP server
	go func() {
		log.Print(WebServerName + " Listening on HTTP port:" +WebServerHTTPPort+ "...")
		err := http.ListenAndServe(":"+WebServerHTTPPort, nil)
		if err != nil {
			log.Fatal("Error starting HTTP server: ", err)
		}
	}()

	// HTTPS server
	go func() {
		HttpsPort, err := strconv.Atoi(WebServerHTTPPort)
		HttpsPort += 300
		log.Print(WebServerName + " Listening on HTTPS port:" + strconv.Itoa(HttpsPort) + "...")
		certFile := "./publiccert.crt"
		keyFile := "./privatekey.key"
		err = http.ListenAndServeTLS(":"+strconv.Itoa(HttpsPort), certFile, keyFile, nil)
		if err != nil {
			log.Fatal("Error starting HTTPS server: ", err)
		}
	}()
	select {}
}

// --------------------------------------------------------------------------------------
//
//	Add CORS headers to HTTP responses
//
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
	// get port number I am listening on
	addOriginHeader = true
	param1 := r.URL.Query().Get("corssettings")
	switch param1 {
	case "TurnCorsOff":
		addOriginHeader = false
	case "TurnCorsWildOn":
		AllowOrigin = "*"
	case "TurnCorsRandomOrigOn": // random port and https
		AllowOrigin = "https://xyz.com:123"
	case "TurnCorsSelfOrigOn":
		AllowOrigin = "http://" + r.Host
	default:
		AllowOrigin = "*"
		addOriginHeader = true
	}
	WriteACHeader(w, AllowOrigin)
	http.Error(w, "Return to Main Page", http.StatusNoContent)
}

// --------------------------------------------------------------------------------------
//      Login Logic 
//      Pretend to login and then return access token in cookie
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
	tokenval := "blank"
	if success {
		// when successful login, create access token in cookie
		tokenname := "AccessToken" // return fake access token
		tokenval = "12345678990"
		// get origin domain to put into cookie
		Origin, _ := url.Parse(r.Header["Origin"][0]) // Origin= http://localhost:9081
		Domain, _, _ := net.SplitHostPort(Origin.Host) // Split Origin
		// create cookie to return access token
		var cookie http.Cookie
		cookie = http.Cookie{Name: tokenname, Value: tokenval, Domain: Domain, Secure: false, SameSite: http.SameSiteLaxMode}	
		http.SetCookie(w, &cookie)
		message = "Login successful!, returning cookie with access token in lax, non-Secure mode"
	} else {
		message = "Invalid username or password, no access token or cookie"
	}

	loginResp := LoginResponse{
		Message: message,
		Token:   tokenval,
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
		w.WriteHeader(http.StatusNoContent)
		WriteACHeader(w, AllowOrigin)
		r.ParseForm()
		// logic part of log in
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}
///////////////////////////////////////////////////////////////////////////
// Cookie Handler
// Add various cookies to response to see if client can read them
///////////////////////////////////////////////////////////////////////////
func MakeVariousCookies(r *http.Request, w http.ResponseWriter){
// Origin and Target should be http://somehostname:portnumber
	var Scheme, Origin, Target string = "","",""
	// What scheme was this called with??
	if (r.TLS == nil) {
		Scheme = "http"
	} else {
		Scheme = "https"
	}
	var Domain string
	// find where the request comes from and what host it is targeting
	 OriginVal,ok := r.Header["Origin"]   // r.Header["Origin"][0]= localhost:9081 
	if (!ok) 	{
		Origin = "NoOriginSpecified" // Oh OH, no origin in Header...
	} else {
		Origin = OriginVal[0]
	}
	// Create Target of my Web server how it was called
	Target = Scheme + "://" + r.Host // r.Host= localhost:9381 
	MyHostArray := strings.Split(r.Host,":")
	Domain = MyHostArray[0]
	// Create various Cookies
	var cookie http.Cookie
	cookievalue := "Origin From: " + Origin + " To: " + Target
	cookie = http.Cookie{Name: "CookieNotSecureLax", Value: cookievalue, Domain: Domain, Secure: false, SameSite: http.SameSiteLaxMode}		
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "CookieSecureLax", Value: cookievalue, Domain: Domain, Secure: true, SameSite: http.SameSiteLaxMode}		
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "CookieNotSecureNone", Value: cookievalue, Domain: Domain, Secure: false, SameSite: http.SameSiteNoneMode}	
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "CookieSecureNone", Value: cookievalue, Domain: Domain, Secure: true, SameSite: http.SameSiteNoneMode}	
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "CookieNotSecureLax_HTTPOnly", Value: cookievalue, Domain: Domain, Secure: false, HttpOnly: true, SameSite : http.SameSiteLaxMode}	
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "CookieNotSecureLax_NOHTTPOnly", Value: cookievalue, Domain: Domain, Secure: false, HttpOnly: false, SameSite : http.SameSiteLaxMode}	
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "CookieNotSecureLax_WeirdDomain", Value: cookievalue, Domain: "WeirdDomain.com", Secure: false, HttpOnly: false, SameSite : http.SameSiteLaxMode}	
	http.SetCookie(w, &cookie)
	return
}


type Message struct {
	Text string `json:"text"`
}

func cookiehandler(w http.ResponseWriter, r *http.Request) {
	// create all sorts of cookies
	MakeVariousCookies(r,w)
	// write the allow-origin header
	WriteACHeader(w, AllowOrigin)
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	// Create an info message
	message := Message{Text: "Set various cookie values"}
	// Convert the message to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Write the JSON data to the response body
	w.Write(jsonData)
}
// --------------------------------------------------------------------------------------
//
//	Add HTTP Request Handler to recieve GET /get-json request to return data to client
//	so see if client can read it cross-origin
//
// --------------------------------------------------------------------------------------


func jsonhandler(w http.ResponseWriter, r *http.Request) {
	// Create a sample message
	message := Message{Text: "ThisPasswordIsSecretFor:" + WebServerName}

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
