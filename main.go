package main

/*
                               CORS Lab

This program is a CORS security lab the allows users to explore both the client and server
side of CORS. Users can manipulate both the Client JavaScript and Server GO HTTP header CORS
attributes and view what the results are.

See README.md for details

*/

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dreezman/browser-security/common"
	"github.com/dreezman/browser-security/cors"
	"github.com/dreezman/browser-security/csrf"
	"github.com/dreezman/browser-security/login"
)

// --------------------------------------------------------------------------------------
//
//	MAIN: is just a web server handling requests
//
// --------------------------------------------------------------------------------------
func main() {
	if len(os.Args) < 2 {
		log.Fatal("Need 2 args, common.WebServerName and PortNumber")

	}
	common.WebServerName = os.Args[1]
	common.WebServerHTTPPort = os.Args[2]

	// Setup all the paths to handle HTTP requests
	fs := http.FileServer(http.Dir("./static"))
	// handle default web requests as front end server for html pages
	http.Handle("/", addHeaders(fs))

	// handle toggle CORS headers
	http.HandleFunc("/cors-toggle", cors.CorsToggle)

	// handle login forms from JS Fetch requests
	http.HandleFunc("/login", login.LoginHandler)
	// handle login forms from classic Form Post Submit
	http.HandleFunc("/classic-form-submit", login.ClassicFormSubmit)

	// test program: return secrets to client, see if they read it
	http.HandleFunc("/get-json", cors.Jsonhandler)

	// set cookies in response
	http.HandleFunc("/get-cookies", csrf.Cookiehandler)

	// HTTP server
	go func() {
		log.Print(common.WebServerName + " Listening on HTTP port:" + common.WebServerHTTPPort+ "...")
		err := http.ListenAndServe(":"+ common.WebServerHTTPPort, nil)
		if err != nil {
			log.Fatal("Error starting HTTP server: ", err)
		}
	}()

	// HTTPS server
	go func() {
		HttpsPort, err := strconv.Atoi(common.WebServerHTTPPort)
		HttpsPort += 300
		log.Print(common.WebServerName + " Listening on HTTPS port:" + strconv.Itoa(HttpsPort) + "...")
		certFile := "./publiccert.crt"
		keyFile := "./privatekey.key"
		err = http.ListenAndServeTLS(":"+strconv.Itoa(HttpsPort), certFile, keyFile, nil)
		if err != nil {
			log.Fatal("Error starting HTTPS server: ", err)
		}
	}()
	select {}
}


func addHeaders(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		common.WriteACHeader(w, common.AllowOrigin)
		fs.ServeHTTP(w, r)
	}

}