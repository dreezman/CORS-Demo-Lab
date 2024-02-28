package csrf

import (
	"browser-security-lab/src/common"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// /////////////////////////////////////////////////////////////////////////
// Cookie Handler
// Add various cookies to response to see if client can read them
// /////////////////////////////////////////////////////////////////////////
func MakeVariousCookies(r *http.Request, w http.ResponseWriter) {
	// Origin and Target should be http://somehostname:portnumber
	var Scheme, Origin, Target string = "", "", ""
	// What scheme was this called with??
	if r.TLS == nil {
		Scheme = "http"
	} else {
		Scheme = "https"
	}
	var Domain string
	// find where the request comes from and what host it is targeting
	OriginVal, ok := r.Header["Origin"] // r.Header["Origin"][0]= localhost:9081
	if !ok {
		Origin = "NoOriginSpecified" // Oh OH, no origin in Header...
	} else {
		Origin = OriginVal[0]
	}
	// Create Target of my Web server how it was called
	Target = Scheme + "://" + r.Host          // r.Host= localhost:9381
	MyHostArray := strings.Split(r.Host, ":") // find the domainname before the "":"
	Domain = MyHostArray[0]
	// Create various Cookies
	var cookie http.Cookie
	cookieOriginTarget := "Origin From: " + Origin + " To: " + Target
	cookie = http.Cookie{Name: "Cookie1", Value: cookieOriginTarget + "||Secure=False||SameSite=Lax||Domain=" + Domain + "|| ", Domain: Domain, Secure: false, SameSite: http.SameSiteLaxMode}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "Cookie2", Value: cookieOriginTarget + "||Secure=True||SameSite=Lax||Domain=" + Domain + "|| ", Domain: Domain, Secure: true, SameSite: http.SameSiteLaxMode}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "Cookie3", Value: cookieOriginTarget + "||Secure=false||SameSite=None||Domain=" + Domain + "|| ", Domain: Domain, Secure: false, SameSite: http.SameSiteNoneMode}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "Cookie4", Value: cookieOriginTarget + "||Secure=true||SameSite=None||Domain=" + Domain + "|| ", Domain: Domain, Secure: true, SameSite: http.SameSiteNoneMode}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "Cookie5", Value: cookieOriginTarget + "||Secure=False||HttpOnly: true||SameSite=Lax||Domain=" + Domain + "|| ", Domain: Domain, Secure: false, HttpOnly: true, SameSite: http.SameSiteLaxMode}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "Cookie6", Value: cookieOriginTarget + "||Secure=False||SameSite=Strict||Domain=" + Domain + "|| ", Domain: Domain, Secure: false, SameSite: http.SameSiteStrictMode}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "Cookie7", Value: cookieOriginTarget + "||Secure=False||SameSite=Lax||Domain=WeirdDomain.com|| ", Domain: "WeirdDomain.com", Secure: false, HttpOnly: false, SameSite: http.SameSiteLaxMode}
	http.SetCookie(w, &cookie)
	return
}

func Cookiehandler(w http.ResponseWriter, r *http.Request) {
	// create all sorts of cookies
	MakeVariousCookies(r, w)
	// write the allow-origin header
	common.WriteACHeader(w, common.AllowOrigin)
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	// Create an info message
	message := common.Message{Text: "Set various cookie values"}
	// Convert the message to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Write the JSON data to the response body
	w.Write(jsonData)
}

// pretend to set a password
func FakeSetPassword(w http.ResponseWriter, r *http.Request) {
	var newpassword string
	switch r.Method {
	// if GET, then get password from URL
	case "GET":
		newpassword = r.URL.Query().Get("new-password")
		fmt.Fprintf(w, "Received GET request to change password to: = %v\n", newpassword)
		// if POST, then get password from body
	case "POST":
		// parse the header into key, value pairs so we can find the password
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Received POST request to change password: r.PostFrom = %v\n", r.PostForm)
		// get new password
		newpassword = r.FormValue("new-password")
	}
	// write the allow-origin header
	common.WriteACHeader(w, common.AllowOrigin)
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	// Create an info message
	message := common.Message{Text: "Setting your password to:" + newpassword}
	// Convert the message to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Write the JSON data to the response body
	w.Write(jsonData)
}
