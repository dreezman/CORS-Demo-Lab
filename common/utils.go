package common

import (
	"net/http"
	"strconv"
)

// Common variables for configuring headers
var AddOriginHeader = true   // add Access-Control header to HTTP response
var AddCredsHeader = false   // add Access-Control header to send credentials
var AllowOrigin string = "*" // Choose a Access-Control origin header, default is allow cross origin all
// Common variables Webserver info
var WebServerName string = ""
var WebServerHTTPPort string = "80"
var WebServerHTTPSPort_Increment = 300 // web server https port is http port + 300

func WriteACHeader(w http.ResponseWriter, AllowOrigin string) {
	if AddOriginHeader {
		//w.Header().Add("X-Frame-Options", "GOFORIT")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Access-Control-Allow-Origin", AllowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Headers: X-PINGOTHER, Content-Type,Cache-Control, Content-Length,Content-Type,Expires,Last-Modified")
	}
	if AddCredsHeader {
		w.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(AddCredsHeader))
	}
}


type Message struct {
	Text string `json:"text"`
}