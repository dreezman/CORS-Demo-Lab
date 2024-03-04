package main

/*
                               Browser Security Lab

This program is a CORS security lab the allows users to explore both the client and server
side of CORS. Users can manipulate both the Client JavaScript and Server GO HTTP header CORS
attributes and view what the results are.

See README.md for details

*/

import (
	"browser-security-lab/src/common"
	"browser-security-lab/src/cors"
	"browser-security-lab/src/csp"
	"browser-security-lab/src/csrf"
	"browser-security-lab/src/login"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/handlers"
)

// Set HTTP request handlers
func handleRequest(frameName string, port string, mux *http.ServeMux) {

	// Setup all the paths to handle HTTP requests
	fs := http.FileServer(http.Dir("./static"))

	// handle default web requests as front end server for html pages
	mux.Handle("/", addHeaders(fs))

	// handle toggle CORS headers
	mux.HandleFunc("/cors-toggle", cors.CorsToggle)

	// handle login forms from JS Fetch requests
	mux.HandleFunc("/login", login.LoginHandler)
	// handle login forms from classic Form Post Submit
	mux.HandleFunc("/classic-form-submit", login.ClassicFormSubmit)
	// fake setting new password
	mux.HandleFunc("/change-password", csrf.FakeSetPassword)
	// test program: return secrets to client, see if they read it
	mux.HandleFunc("/get-json", cors.Jsonhandler)

	// set cookies in response
	mux.HandleFunc("/get-cookies", csrf.Cookiehandler)
	// set cookies in response

	mux.HandleFunc("/xss-attack", cors.XssAttackHandler)
	// set CSP Header global vars
	mux.HandleFunc("/set-csp-header", csp.SetCSPHeader)
	mux.HandleFunc("/write-csp-header", csp.WriteCSPHeader)
	// Handle and print all CSP violations
	mux.HandleFunc("/csp-report-only", csp.CSPReportOnlyHandler)

}

// Add Cors headers to all responses
func addHeaders(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		common.WriteACHeader(w, common.AllowOrigin)
		fs.ServeHTTP(w, r)
	}

}

// --------------------------------------------------------------------------------------
//
//	MAIN: is just a web server handling requests
//
// --------------------------------------------------------------------------------------
func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Hey you forgot to tell me where the config file is located!! Missing 1 arg: path to iframe config file ")

	}

	// read iframe config file
	jsonFile, err := os.Open(os.Args[1])
	if err != nil || jsonFile == nil {
		fmt.Println("Oh Oh..cannot read frame config file, do you have the right path??: ", err)
		return
	}
	// load JSON into GO object
	common.LoadFrameConfig(jsonFile)
	if err != nil {
		fmt.Println("Something went wrong when loading frame config file, is the data corrupt or not JSON??: ", err)
		return
	}

	// Each Iframe will be served by a mini-backend web server.
	// Each mini-webserver will be a background thread/goroutine dedicated
	// to processing those iframe requests.
	//
	// Loop thru all the iframes in the config file and
	// setup a background goroutine with an HTTP server for each iframe
	// and a HTTPS server for each iframe. So 3 iframes, 6 web servers
	//
	// NOTE: to setup HTTPS
	// HTTPS server
	// import to certmgr->trusted certs
	// openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 \
	//  -nodes -keyout privatekey.key -out publiccert.crt -subj "/CN=localhost" \
	//  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
	// firefox: add exception
	//  https://unix.stackexchange.com/questions/644176/how-to-permanently-add-self-signed-certificate-in-firefox

	var wg sync.WaitGroup
	for frameName, frameData := range common.IFrameConfigMap {
		// HTTP server
		wg.Add(1)
		go func(frameName string, frameData common.Frame) {
			defer wg.Done()
			mux := http.NewServeMux()
			handleRequest(frameName, frameData.HTTPPort, mux)
			log.Print(frameName + " Listening on HTTP port:" + frameData.HTTPPort + "...")
			err := http.ListenAndServe(":"+frameData.HTTPPort, handlers.LoggingHandler(os.Stdout, mux))
			if err != nil {
				log.Fatal("Error starting web server: ", err)
			}
		}(frameName, frameData)
		// HTTPS server
		wg.Add(1)
		go func(frameName string, frameData common.Frame) {
			defer wg.Done()
			mux := http.NewServeMux()
			handleRequest(frameName, frameData.HTTPSPort, mux)
			log.Print(frameName + " Listening on HTTPS port:" + frameData.HTTPSPort + "...")
			certFile := "./publiccert.crt"
			keyFile := "./privatekey.key"
			err = http.ListenAndServeTLS(":"+frameData.HTTPSPort, certFile, keyFile, handlers.LoggingHandler(os.Stdout, mux))
			if err != nil {
				log.Fatal("Error starting web server: ", err)
			}
		}(frameName, frameData)
	}
	wg.Wait()

}
