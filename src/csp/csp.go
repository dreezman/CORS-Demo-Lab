package csp

import (
	"browser-security-lab/src/common"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

// CSP Data structure
// to create a CSP header that looks like this
// Content-Security-Policy: default-src 'self' example.com *.example.com; object-src 'none'; base-uri 'self';
type CSP_Data struct {
	CSP_Type string   `json:"csp-type"` // default-src, script-src, etc
	Domains  []string `json:"domains"`
}

type CSP_Struct struct {
	Enabled  bool       `json:"enabled"`  // if nothing is enabled, then CSP is disabled
	CSP_Mode string     `json:"cspMode"`  // Content-Security-Policy, Content-Security-Policy-Report-Only
	CSP_Data []CSP_Data `json:"csp-data"` // default-src, script-src, etc
}

var CSPConfig_Current CSP_Struct = CSP_Struct{Enabled: false} // default is to disable csp
var CSPHeader string = ""                                     // Content-Security-Policy, Content-Security-Policy-Report-Only
var CSPDomains string = ""                                    // default-src 'self'; img-src *; media-src example.org example.net; script-src userscripts.example.com

// Set the CSP header based on the input, do not write the header to the response
func SetCSPHeader(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&CSPConfig_Current)
	if err != nil {
		panic(err)
	}
	log.Println(CSPConfig_Current)

	CSPHeader = ""
	CSPDomains = ""
	if CSPConfig_Current.Enabled {

		if CSPConfig_Current.CSP_Mode == "Content-Security-Policy" {
			CSPHeader = "Content-Security-Policy"
		} else if CSPConfig_Current.CSP_Mode == "Content-Security-Policy-Report-Only" {
			CSPHeader = "Content-Security-Policy-Report-Only"
		} else {
			http.Error(w, "Invalid CSP mode: "+CSPConfig_Current.CSP_Mode, http.StatusBadRequest)
			CSPConfig_Current.Enabled = false
			return
		}

		if len(CSPConfig_Current.CSP_Data) == 0 {
			CSPDomains = "default-src 'self'; "
		}
		// build the header to look like this
		// Content-Security-Policy: default-src 'self' example.com *.example.com; object-src 'none'; base-uri 'self';
		for _, cspData := range CSPConfig_Current.CSP_Data {
			CSPDomains += cspData.CSP_Type + " " // default-src, script-src, etc
			for _, domain := range cspData.Domains {
				CSPDomains += domain + " " // abc.com 123.com;
			}
			CSPDomains += ";"
		}
	}
}

func InsertCSPHeader(w http.ResponseWriter, r *http.Request) {

	// Add dynamic headers based on request properties
	if CSPConfig_Current.Enabled {
		url := common.IFrameConfigMap["ParentIframe"].FullHTTPSURL + "/csp-report-only"
		cspGroup := `{"group": "csp-endpoint-group","max_age": 10886400,"endpoints": [{"url": "` + url + `" }]}`
		w.Header().Set("Report-To", cspGroup)
		w.Header().Set("Reporting-Endpoints", `csp-endpoint-uri="`+url+`"`)
		w.Header().Set(CSPHeader, CSPDomains) // CSP Header
	}

}

/*
// Write the CSP header into the HTTP response
func InsertCSPHeader() map[string]string {

	headers := make(map[string]string)
	// Add dynamic headers based on request properties
	if CSPConfig_Current.Enabled {
		url := common.IFrameConfigMap["ParentIframe"].FullHTTPSURL + "/csp-report-only"
		print(fmt.Sprint("csp-endpoint-uri=\"", url, "\""))
		cspGroup := `{"group": "csp-endpoint-group","max_age": 10886400,"endpoints": [{"url": "` + url + `" }]}`
		headers["Report-To"] = cspGroup
		headers["Reporting-Endpoints"] = `csp-endpoint-uri="` + url + `"`
		headers[CSPHeader] = CSPDomains // CSP Header
	}

	return headers
}
*/
// --------------------------------------------------------------------------------------
// Handle all CSP violations and print them out
// --------------------------------------------------------------------------------------
type CSPReport struct {
	CSPReport CSP `json:"csp-report"`
}

type CSP struct {
	DocumentURI string `json:"document-uri"`
	Referrer    string `json:"referrer"`
	BlockedURI  string `json:"blocked-uri"`
}

func CSPReportOnlyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	//	var report CSPReport
	//	err = json.Unmarshal(body, &report)
	var report interface{}
	err = json.Unmarshal([]byte(body), &report)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshalling JSON csp report: %v\n", err)
		fmt.Fprintf(os.Stderr, "CSP Body: %s\n", string(body))
		return
	}

	//reportJson, err := json.Marshal(report)
	reportJson, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling csp report: %v", err)
		fmt.Fprintf(os.Stderr, "CSP Body: %s\n", string(body))
		return
	}

	fmt.Fprintf(os.Stderr, "Content-Security-Policy-Report-Only: %s\n", reportJson)
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

// --------------------------------------------------------------------------------------
// Add HTTP Request Handler to recieve GET /xss-attack request to return data to client
// that is a XSS string
// --------------------------------------------------------------------------------------
func XssFormHandler(w http.ResponseWriter, r *http.Request) {

	fname := r.URL.Query().Get("fname")
	//lname := r.URL.Query().Get("fname")
	//fmt.Fprintf(w, "Received GET XSS request with XSS as value: = %v\n", xssVal)
	// Write the XSS data to the response body
	log.Print("XSS attack response: ", fname)
	w.Write([]byte(fname))
}
