package common

import (
	"encoding/json"
	"fmt"
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

// Push the CSP config into the Nginx config file
func PushNgxConfig(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is OPTIONS
	if r.Method == http.MethodOptions {
		// Write the Access Control header to allow all requests from all origins
		AllowOrigin = "*"
		fmt.Println("in options")
		WriteACHeader(w, AllowOrigin)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Anyone allowed to POST"))
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// Write the Access Control header to allow all requests from all origins
	AllowOrigin = "*"
	WriteACHeader(w, AllowOrigin)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	bodyStr := string(body)
	fmt.Print("Before file exists check, here is body:\n", bodyStr)
	// Check if the file exists
	filePath := "/usr/share/nginx-config/ngx-dynamic-update-file.conf"
	var ngxConfigFile *os.File // Declare ngxConfigFile before the conditional statements
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		w.WriteHeader(http.StatusAccepted) // 202 Accepted
		w.Write([]byte(fmt.Sprintf("Warning! NGX ConfigFile %s does not exist yet, nothing POSTed\n", filePath)))
		return
	} else if err != nil { // error checking if file exists
		http.Error(w, fmt.Sprintf("Error checking if file %s exists", filePath), http.StatusInternalServerError)
		return
	} else {
		// kill file contents
		err := os.Truncate(filePath, 0)
		if err != nil {
			fmt.Print("Truncating ngx.conf file error: \n", err)
			http.Error(w, "Truncating ngx.conf file error:", http.StatusInternalServerError)
			return;
		}
		// Open the file
		fmt.Print("Before file open\n")
		ngxConfigFile, err = os.OpenFile(filePath, os.O_RDWR, 0666)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error opening Nginx config file %s for writing", filePath), http.StatusInternalServerError)
			return
		}
		defer ngxConfigFile.Close()
	}

	// Write the body into the file
	// Convert the string to ASCII
	asciiBody := ""
	for _, char := range bodyStr {
		if char > 127 {
			asciiBody += "?" // Replace non-ASCII characters with a placeholder
		} else {
			asciiBody += string(char)
		}
	}

	// Write the body into the file
	fmt.Print("Before write, here is body: \n", asciiBody,"\n")
	_, err = ngxConfigFile.Write([]byte(asciiBody))
	if err != nil {
		fmt.Print("Write error: \n", err)
		http.Error(w, "Error writing to Nginx config file", http.StatusInternalServerError)
		return
	}
	// Write the body into the file
	fmt.Print("\nSnyc file to write: \n")
	err = ngxConfigFile.Sync()
	if err != nil {
		fmt.Print("Snyc error: \n", err)
		http.Error(w, "Error syncing to Nginx config file", http.StatusInternalServerError)
		return
	}
	fmt.Print("Before write cors\n")
	w.WriteHeader(http.StatusOK)
	AllowOrigin = "*"
	WriteACHeader(w, AllowOrigin)
	w.Write([]byte("Nginx config updated successfully"))
	ngxConfigFile.Close()
}
