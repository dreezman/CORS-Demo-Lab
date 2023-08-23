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
	"io"
	"log"
	"net/http"
	"os"

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
	if len(os.Args) < 1 {
		log.Fatal("Missing 1 arg: path to iframe config file ")

	}

	// read iframe config file
	jsonFile, err := os.Open(os.Args[0])
    if (err != nil || jsonFile == nil ){
       fmt.Println(err)
	}
	byteValue, _ := io.ReadAll(jsonFile)  /// stuck here reading in map	
	err = json.Unmarshal([]byte(byteValue), &common.FrameConfigData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
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
	// fake setting new password
	http.HandleFunc("/change-password",csrf.FakeSetPassword)
	// test program: return secrets to client, see if they read it
	http.HandleFunc("/get-json", cors.Jsonhandler)

	// set cookies in response
	http.HandleFunc("/get-cookies", csrf.Cookiehandler)


	for frameName, frameData := range common.FrameConfigData {
		// HTTP server
		go func() {
			log.Print(frameName + " Listening on HTTP port:" + frameData.HttpPort + "...")
			err := http.ListenAndServe(":" + frameData.HttpPort, nil)
			if err != nil {
				log.Fatal("Error starting HTTP server: ", err)
			}
		}()

		// HTTPS server
			// import to certmgr->trusted certs
			// openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 \
			//  -nodes -keyout privatekey.key -out publiccert.crt -subj "/CN=localhost" \
			//  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
			// firefox: add exception
			//  https://unix.stackexchange.com/questions/644176/how-to-permanently-add-self-signed-certificate-in-firefox		

		go func() {
			log.Print(frameName + " Listening on HTTPS port:" + frameData.HttpsPort + "...")
			certFile := "./publiccert.crt"
			keyFile := "./privatekey.key"
			err = http.ListenAndServeTLS(":"+frameData.HttpsPort, certFile, keyFile, nil)
			if err != nil {
				log.Fatal("Error starting HTTPS server: ", err)
			}
	}()}
	select {}
}


func addHeaders(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		common.WriteACHeader(w, common.AllowOrigin)
		fs.ServeHTTP(w, r)
	}

}