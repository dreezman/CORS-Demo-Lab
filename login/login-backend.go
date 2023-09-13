package login

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/dreezman/browser-security/common"
)

// --------------------------------------------------------------------------------------
//      Login Logic
//      Pretend to login and then return access token in cookie
// --------------------------------------------------------------------------------------

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	Success bool   `json:"success"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		common.WriteACHeader(w, common.AllowOrigin)
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Perform authentication logic here

	// Check if the username and password are valid
	success := loginReq.Username == "admin" && loginReq.Password == "password"
	message := ""
	tokenval := "blank"
	if success {
		// when successful login, create access token in cookie
		tokenval = "12345678990"
		// get origin domain to put into cookie
		Origin, _ := url.Parse(r.Header["Origin"][0]) // Origin= http://localhost:9081
		Domain, _, _ := net.SplitHostPort(Origin.Host) // Split Origin
		// create cookie to return several access tokens, just to see what happens
		var cookie http.Cookie

		tokenname := "AccessToken_LaxSameSite_NotSecure"; // return fake access token
		cookie = http.Cookie{Name: tokenname, Value: tokenval, Domain: Domain, Secure: false, SameSite: http.SameSiteLaxMode}	
		http.SetCookie(w, &cookie)

		tokenname  = "AccessToken_NoSameSite_NotSecure"; // return fake access token
		cookie = http.Cookie{Name: tokenname, Value: tokenval, Domain: Domain, Secure: false, SameSite: http.SameSiteNoneMode}	
		http.SetCookie(w, &cookie)

		tokenname  = "AccessToken_NoSameSite_Secure"; // return fake access token
		cookie = http.Cookie{Name: tokenname, Value: tokenval, Domain: Domain, Secure: true, SameSite: http.SameSiteNoneMode}	
		http.SetCookie(w, &cookie)

		message = "Login successful!, returning cookie with access token in lax, non-Secure mode"
	} else {
		message = "Invalid username or password, no access token or cookie"
	}

	loginResp := LoginResponse{
		Message: message,
		Token:   tokenval,
		Success: success,
	}
	w.Header().Set("Content-Type", "application/json")
	common.WriteACHeader(w, common.AllowOrigin)
	if err := json.NewEncoder(w).Encode(loginResp); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}

func ClassicFormSubmit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		common.WriteACHeader(w, common.AllowOrigin)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Must use a POST, preferabally HTTPS"))
	} else {
		w.WriteHeader(http.StatusNoContent)
		common.WriteACHeader(w, common.AllowOrigin)
		r.ParseForm()
		// logic part of log in
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}