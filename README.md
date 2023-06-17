- [CORS Demo Lab](#cors-demo-lab)
  - [Installation](#installation)
  - [Usage](#usage)
    - [Start Program](#start-program)
    - [Cross-Origin Queries](#cross-origin-queries)
  - [Contributing](#contributing)
  - [License](#license)



# CORS Demo Lab

One of the elements missing when watching YouTube videos explaining technology, is the hands on experience
to explore the "What Ifs" along with seeing both sides of the client-server interaction on a step-by-step
basis. This is the only way to fully understand how a technology works.


Same-Origin-Policy (SOP) and Cross-Origin-Requests (CORS) are easy to understand at a high
level, but once one reads all the rules behind it, becomes very complex. There are a lot of
"Rule XXX is always true EXCEPT in these cases". So in order to explore the intracies of 
CORS "What-Ifs", I created a Golang based CORS lab.

Lesson 1: Site vs Origin..
Same Origin Policy states that JavaScript cannot access data across Origins. What is an Origin?
An origin is very strict, the whole url domain from 'h' to portnumber will make a match. Subdomains are not in the same Origin.
Site is more flexible, subdomains are in the same Site. These are used in cookie restrictions.

![Alt text](images/origin-site.jpg)

This CORS security lab allows users to explore both the client and server
side of CORS. Users can manipulate both the Client JavaScript and Server GO HTTP header CORS
attributes and view what the results are.

This is done thru the use of a main program and two iframes, all in different origins via
unique port numbers. The main program is a web server that forks to sub-processes all with
different port numbers which makes the document.location.origin unique. 

![Alt text](images/cors-lab.jpg)

Users can then switch between the 3 JS contexts (main, iframe1, iframe2) and view how CORS impacts accessing data
from the different origins. Users can do this with:

- Forms - Post login and try to view response
- getElementById() to retrieve data between origins and see response
- postMessage between iframes to retrieve data
- JS fetch - can one fetch/HTTP GET and view response between origins
- localStorage - Can one access local storage between origins



## Installation

—-------------
 Create a golang HELLO WORLD project---------------
- https://go.dev/doc/install
- Download     https://go.dev/dl/
- https://go.dev/doc/tutorial/getting-started

—-------------
 Install CORS lab ---------------
- cd ..
- git clone https://github.com/dreezman/CORS-Demo-Lab.git

```
// Unix
kill $(jobs -p) ; sleep 3 ; go run main.go TLD 8081 & go run main.go iframe1 3000 & go run main.go iframe2 3001 &

// Windows
get-job| stop-job | remove-job ; go run main.go TLD 8081 & go run main.go iframe1 3000 & go run main.go iframe2 3001 &
```

## Usage

We will be exploring:

- cross-domain inspections internally inside browser
- postMessage (inter-window messaging) vs fetch(HTTP requests)

We will use the browser debugger/Inspector a lot to examine the code and traffic between the client and browser. 
- Chrome:   Press F12 in Chrome or 
- FireFox: 3-bars-MoreTools-WebDeveloperTools

to begin inspection.
<hr>

### Start Program

Wait for program to start, takes 30 seconds, then browse to main web server it will load the HTML into all the 3 web servers from each
or their respective origins. 
```
http://localhost:8081/?iframeurl1=http://localhost:3000/iframes.html&iframeurl2=http://localhost:3001/iframes.html
```

You can see pages loaded in their respective origins.

![Alt text](images/iframe-setup.jpg)


Now press on the "Send Message" buttons to send messages postMessages between the 
iframes and the parent, all in different origins. 

Use Chrome Inspector or Firefox Debugger, to monitor the network traffic and inspect the header fields for the CORS headers. 
![Alt text](images/cors-headers.jpg)

The default is to allow all cross-origin requests via "Access-Control-Allow-Origin: *"

You can then modify the [main.go](./main.go) program to change HTTP header fields and watch the CORS errors occur in the console.
```
const addOriginHeader = true // add Access-Control header to HTTP response
var AllowOrigin string = "*" // Choose a Access-Control origin header
//var AllowOrigin string = "http://localhost:8081"
//var AllowOrigin string = "http://localhost:3000"
//var AllowOrigin string = "http://localhost:3001"
//var AllowOrigin string = "http://localhost:222"
```
<hr>

Also test the login sequence with user: admin password: password and view the post and response. Note the response token is stored in the localStorage.
![Alt text](images/login.jpg)
![Alt text](images/token.jpg)





### Cross-Origin Queries

Now we will manually to cross-Origin queries by manually executing JavaScript commands.

In the Inspector console, one can switch between JavaScript contexts to execute JavaScript commands in the local context.

![Alt text](images/jscontexts.jpg)


Use the Inspector Console to manually do queries between the 3 origins to see if you can read the responses to the queries.


```
response=await fetch("http://localhost:8081/get-json"); await response.text()
response=await fetch("http://localhost:3000/get-json"); await response.text()
```
![Alt text](images/fetch-queries.jpg) 
<hr>

Also try reading local storage to/from different origins to see if you can access local storage.


```
# From Parent to Iframes
window.frames[0].localStorage
window.frames[1].localStorage
# From Iframes to each other
window.parent.frames[0].localStorage
window.parent.frames[1].localStorage
# From Iframes to parent
window.parent.localStorage
```


Video coming soon...


## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

- [MIT](https://choosealicense.com/licenses/mit/)
- [GO](https://go.dev/LICENSE)