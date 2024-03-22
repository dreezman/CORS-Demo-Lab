package common

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
)

// Common variables for configuring headers
var AddOriginHeader = true   // add Access-Control header to HTTP response
var AddCredsHeader = false   // add Access-Control header to send credentials
var AllowOrigin string = "*" // Choose a Access-Control origin header, default is allow cross origin all

// Write the Access Control CORS header into the HTTP response
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

/*
	 -----------------------------------------------------------------------

		Load iFrame config data from file into this structure as a map
		So that at the end we can loop thru the map and refer to each
		frame property like

		FrameConfigMap[frameName].DomainName

		Bit weird, first read JSON config data, the put into map where we
		can refer to each member by name.

	   -----------------------------------------------------------------------
*/
type Frame struct {
	FrameName    string `json:"frameName"`
	DomainName   string `json:"domainName"`
	HTTPPort     string `json:"httpPort"`
	HTTPSPort    string `json:"httpsPort"`
	FullHTTPURL  string `json:"fullHTTPURL"`
	FullHTTPSURL string `json:"fullHTTPSURL"`
	Description  string `json:"Description"`
}

var IFrameConfigMap map[string]Frame

/* First load temp JSON data into this array */
type IframesData struct {
	Iframes []Frame `json:"Iframes"`
}

var frameConfigData IframesData

// Function to read JSON data and then load into final Iframe map
func LoadFrameConfig(configFile *os.File) error {
	byteValue, _ := io.ReadAll(configFile)
	// read in JSON data
	err := json.Unmarshal([]byte(byteValue), &frameConfigData)
	if err != nil {
		return err
	}
	// Build map with framename as index so that can refer to
	// as FrameConfigMap[frameName].DomainName
	IFrameConfigMap = make(map[string]Frame)
	for _, iframe := range frameConfigData.Iframes {
		IFrameConfigMap[iframe.FrameName] = iframe
	}
	return nil
}

type Message struct {
	Text string `json:"text"`
}
