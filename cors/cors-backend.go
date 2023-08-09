package cors

import (
	"encoding/json"
	"net/http"

	"github.com/dreezman/browser-security/common"
)

// --------------------------------------------------------------------------------------
//
//	Add CORS headers to HTTP responses
//
// --------------------------------------------------------------------------------------

// Let browser set the CORS header setting
func CorsToggle(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// What should the common.AllowOrigin be set to, default is Allow CORS
	common.AddOriginHeader = true
	param1 := r.URL.Query().Get("common.AllowOrigin")
	switch param1 {
	case "TurnCorsOff":
		common.AddOriginHeader = false
	case "TurnCorsWildOn":
		common.AllowOrigin = "*"
	case "TurnCorsRandomOrigOn": // random port and https
		common.AllowOrigin = "https://xyz.com:123"
	case "TurnCorsSelfOrigOn":
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
}

// --------------------------------------------------------------------------------------
//
//	Add HTTP Request Handler to recieve GET /get-json request to return data to client
//	so see if client can read it cross-origin
//
// --------------------------------------------------------------------------------------


func Jsonhandler(w http.ResponseWriter, r *http.Request) {
	// Create a sample message
	message := common.Message{Text: "ThisPasswordIsSecretFor:" + common.WebServerName}

	// Convert the message to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	common.WriteACHeader(w, common.AllowOrigin)

	// Write the JSON data to the response body
	w.Write(jsonData)
}