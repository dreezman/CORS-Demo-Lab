package csrf

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dreezman/browser-security/common"
)

///////////////////////////////////////////////////////////////////////////
// Cookie Handler
// Add various cookies to response to see if client can read them
///////////////////////////////////////////////////////////////////////////
func MakeVariousCookies(r *http.Request, w http.ResponseWriter){
// Origin and Target should be http://somehostname:portnumber
	var Scheme, Origin, Target string = "","",""
	// What scheme was this called with??
	if (r.TLS == nil) {
		Scheme = "http"
	} else {
		Scheme = "https"
	}
	var Domain string
	// find where the request comes from and what host it is targeting
	 OriginVal,ok := r.Header["Origin"]   // r.Header["Origin"][0]= localhost:9081 
	if (!ok) 	{
		Origin = "NoOriginSpecified" // Oh OH, no origin in Header...
	} else {
		Origin = OriginVal[0]
	}
	// Create Target of my Web server how it was called
	Target = Scheme + "://" + r.Host // r.Host= localhost:9381 
	MyHostArray := strings.Split(r.Host,":")
	Domain = MyHostArray[0]
	// Create various Cookies
	var cookie http.Cookie
	cookievalue := "Origin From: " + Origin + " To: " + Target
	cookie = http.Cookie{Name: "CookieNotSecureLax", Value: cookievalue, Domain: Domain, Secure: false, SameSite: http.SameSiteLaxMode}		
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "CookieSecureLax", Value: cookievalue, Domain: Domain, Secure: true, SameSite: http.SameSiteLaxMode}		
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "CookieNotSecureNone", Value: cookievalue, Domain: Domain, Secure: false, SameSite: http.SameSiteNoneMode}	
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "CookieSecureNone", Value: cookievalue, Domain: Domain, Secure: true, SameSite: http.SameSiteNoneMode}	
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "CookieNotSecureLax_HTTPOnly", Value: cookievalue, Domain: Domain, Secure: false, HttpOnly: true, SameSite : http.SameSiteLaxMode}	
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "CookieNotSecureLax_NOHTTPOnly", Value: cookievalue, Domain: Domain, Secure: false, HttpOnly: false, SameSite : http.SameSiteLaxMode}	
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "CookieNotSecureLax_WeirdDomain", Value: cookievalue, Domain: "WeirdDomain.com", Secure: false, HttpOnly: false, SameSite : http.SameSiteLaxMode}	
	http.SetCookie(w, &cookie)
	return
}




func Cookiehandler(w http.ResponseWriter, r *http.Request) {
	// create all sorts of cookies
	MakeVariousCookies(r,w)
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


func FakeSetPassword(w http.ResponseWriter, r *http.Request) {

	// fake setting password
	// get password from form POST json data
	if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
	}
	fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
	// get new password
	newpassword := r.FormValue("new-password")
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