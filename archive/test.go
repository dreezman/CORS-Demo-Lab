package main

import (
	"fmt"
	"net/http"
	"sync"
)

func startServer(port int, wg *sync.WaitGroup) {
	defer wg.Done()

	addr := fmt.Sprintf(":%d", port)
	fileServer := http.FileServer(http.Dir("static")) // Replace "static" with your actual directory
    mux := http.NewServeMux()
	mux.Handle("/", fileServer)

	fmt.Printf("Server started on port %d\n", port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("Error on port %d: %v\n", port, err)
	}
}

func main() {
	ports := []int{8080, 8081, 8082} // Define the ports you want to listen on
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)
		go startServer(port, &wg)
	}

	wg.Wait()
}