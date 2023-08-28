package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	jsonData := `
	{
		"ParentIframe": {
			"domainName": "localhost",
			"httpPort": "9081",
			"httpsPort": "9381",
			"fullHTTPURL": "http://localhost:9081",
			"fullHTTPSURL": "https://localhost:9381",
			"Description": "ParentFrame has port 9081 and TLS 9381"
		},
		"Iframe1": {
			"domainName": "localhost",
			"httpPort": "3000",
			"httpsPort": "3300",
			"fullHTTPURL": "http://localhost:3000",
			"fullHTTPSURL": "https://localhost:3300",
			"Description": "Iframe1 has different origin as parent, port 3000 and TLS 3300"
		},
		"Iframe2": {
			"domainName": "localhost",
			"httpPort": "3100",
			"httpsPort": "3400",
			"fullHTTPURL": "http://localhost:3001",
			"fullHTTPSURL": "https://localhost:3301",
			"Description": "Iframe2 has different origin as parent, port 3001 and TLS 3301"
		},
		"Iframe3": {
			"domainName": "localhost",
			"httpPort": "9081",
			"httpsPort": "9381",
			"fullHTTPURL": "http://localhost:9081/",
			"fullHTTPSURL": "https://localhost:9381/",
			"Description": "Iframe3 has same origin as Parent and same ports"
		}
	}`

	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print out the contents of the data variable
	for key, value := range data {
		fmt.Println("Key:", key)
		fmt.Println("Value:", value)
		fmt.Println("-------------------------")
	}
}
