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

// CustomHandler is a struct for customizing the file server
// CustomHandler is a struct for customizing the file server
type CustomHandler struct {
	// The base file server
	FileServer http.Handler
	// Function to generate custom headers dynamically
	HeaderFunc func(r *http.Request) map[string]string
}

// ServeHTTP serves the static file and inserts headers
func (h *CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get custom headers from the provided function
	headers := h.HeaderFunc(r)

	// Add custom headers
	for key, value := range headers {
		w.Header().Set(key, value)
	}

	// Serve the static file
	h.FileServer.ServeHTTP(w, r)
}

// Example function to generate headers dynamically
func dynamicHeaders(r *http.Request) map[string]string {
	headers := make(map[string]string)

	// Add dynamic headers based on request properties
	if csp.CSPConfig_Current.Enabled {
		//headers = csp.InsertCSPHeader()
	}
	headers["X-Custom-Header-mike"] = r.URL.Path // Example dynamic header

	return headers
}

func customHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set custom headers here
		csp.InsertCSPHeader(w, r)
		w.Header().Set("X-Custom-Header", "Custom Value")
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Set HTTP request handlers
func handleRequest(mux *http.ServeMux) {

	// Create a file server handler to serve static files from the "static" directory
	fs := http.FileServer(http.Dir("static/"))

	// Define a function to generate custom headers dynamically
	headerFunc := dynamicHeaders

	// Wrap the file server with the custom handler
	customHandler := &CustomHandler{
		FileServer: fs,
		HeaderFunc: headerFunc,
	}
	// Serve static files with the custom handler
	mux.Handle("/", customHandler)
	//fs1 := http.FileServer(http.Dir("/"))
	//mux.Handle("/", fs1)
	/* Specify the directory you want to serve files from
	fs := http.FileServer(http.Dir("static"))
	//--------------------------------------------
	// Use the FileServer to serve files from the specified directory
	mux.Handle("/", fs)
	*/
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

	mux.HandleFunc("/xss-attack", csp.XssAttackHandler)
	// set CSP Header global vars
	mux.HandleFunc("/set-csp-header", csp.SetCSPHeader)
	// Handle and print all CSP violations
	mux.HandleFunc("/csp-report-only", csp.CSPReportOnlyHandler)

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
			handleRequest(mux)
			wrappedMux := customHeaderMiddleware(mux)
			log.Print(frameName + " Listening on HTTP port:" + frameData.HTTPPort + "...")
			err := http.ListenAndServe(":"+frameData.HTTPPort, handlers.LoggingHandler(os.Stdout, wrappedMux))
			if err != nil {
				log.Fatal("Error starting web server: ", err)
			}
		}(frameName, frameData)
		// HTTPS server
		wg.Add(1)
		go func(frameName string, frameData common.Frame) {
			defer wg.Done()
			mux := http.NewServeMux()
			handleRequest(mux)
			wrappedMux := customHeaderMiddleware(mux)
			log.Print(frameName + " Listening on HTTPS port:" + frameData.HTTPSPort + "...")
			certFile := "./publiccert.crt"
			keyFile := "./privatekey.key"
			err = http.ListenAndServeTLS(":"+frameData.HTTPSPort, certFile, keyFile, handlers.LoggingHandler(os.Stdout, wrappedMux))
			if err != nil {
				log.Fatal("Error starting web server: ", err)
			}
		}(frameName, frameData)
	}
	wg.Wait()

}
