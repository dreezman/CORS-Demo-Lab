package csp

import (
	"browser-security-lab/src/common"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

/*
connect-src
https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/connect-src

default-src
https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/default-src

script-src

unsafe-inline	- allow inline scripts
unsafe-eval - allow eval() function

script-src-elem
script-src-attr
*/
// Common variables for configuring headers
var AddCSPHeader = true // add Access-Control header to HTTP response

// --------------------------------------------------------------------------------------
//
//	Add csp headers to HTTP responses
//
// --------------------------------------------------------------------------------------

// Let browser set the csp header setting
/*

input looks like one of these

{
  "enable-csp": {
    "default-src": [
      "abc.com",
      "123.com"
    ],
    "script-src": [
      "abc.com",
      "123.com"
    ]
  },
  "report-only-csp": {
    "default-src": [
      "abc.com",
      "123.com"
    ],
    "script-src": [
      "abc.com",
      "123.com"
    ]
  },
  "disable-csp": {}
}

*/
type CSP struct {
	DefaultSrc []string `json:"default-src"`
	ScriptSrc  []string `json:"script-src"`
}

type CSPConfig struct {
	EnableCSP     CSP `json:"enable-csp"`
	ReportOnlyCSP CSP `json:"report-only-csp"`
	DisableCSP    CSP `json:"disable-csp"`
}

func SetCSPHeader(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var CSPConfig CSPConfig
	err := decoder.Decode(&CSPConfig)
	if err != nil {
		panic(err)
	}
	log.Println(CSPConfig)

	/*
		// What should the common.AllowOrigin be set to, default is Allow csp
		common.AddOriginHeader = true
		param1 := r.URL.Query().Get("AllowOrigin")
		switch param1 {
		case "TurncspOff":
			common.AddOriginHeader = false
		case "TurncspWildOn":
			common.AllowOrigin = "*"
		case "TurncspRandomOrigOn": // random port and https
			common.AllowOrigin = "https://xyz.com:123"
		case "TurncspSelfOrigOn":
			common.AllowOrigin = "http://" + r.Host
		default:
			common.AllowOrigin = "*"
			common.AddOriginHeader = true
		}
		// What should the Send Credentials be set to, default is do not send
		common.AddCredsHeader = false
		param1 = r.URL.Query().Get("creds")
		switch param1 {
		case "Off":
			common.AddCredsHeader = false
		case "On":
			common.AddCredsHeader = true
		default:
			common.AddOriginHeader = false
		}

		common.WriteACHeader(w, common.AllowOrigin)
		http.Error(w, "Return to Main Page", http.StatusNoContent)
	*/
}

// Write the Access Control CORS header into the HTTP response
func WriteCSPHeader(w http.ResponseWriter, AllowOrigin string) {
	if common.AddOriginHeader {
		//w.Header().Add("X-Frame-Options", "GOFORIT")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Access-Control-Allow-Origin", AllowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Headers: X-PINGOTHER, Content-Type,Cache-Control, Content-Length,Content-Type,Expires,Last-Modified")
	}
	if common.AddCredsHeader {
		w.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(common.AddCredsHeader))
	}
}

// --------------------------------------------------------------------------------------
// Add HTTP Request Handler to recieve GET /xss-attack request to return data to client
// that is a XSS string
// --------------------------------------------------------------------------------------
func XssAttackHandler(w http.ResponseWriter, r *http.Request) {

	xssVal := r.URL.Query().Get("xssvalue")
	//fmt.Fprintf(w, "Received GET XSS request with XSS as value: = %v\n", xssVal)
	// Write the XSS data to the response body
	log.Print("XSS attack response: ", xssVal)
	w.Write([]byte(xssVal))
}
