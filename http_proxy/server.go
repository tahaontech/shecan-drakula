package http_proxy

import (
	"fmt"
	"io"
	"net/http"
)

type ProxyServer struct {
	proxyAddr string
}

func NewProxyServer(proxyAddr string) *ProxyServer {
	return &ProxyServer{
		proxyAddr: proxyAddr,
	}
}

func (s *ProxyServer) handleProxy(w http.ResponseWriter, r *http.Request) {
	// Create a new HTTP client
	client := &http.Client{}

	// Forward the request to the target server
	resp, err := client.Do(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending request: %s", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response headers to the client
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set the status code
	w.WriteHeader(resp.StatusCode)

	// Copy the response body to the client
	io.Copy(w, resp.Body)
}

func (s *ProxyServer) Start() {
	// Define the proxy server address and port
	proxyAddr := ":8080"

	// Start the HTTP proxy server
	http.HandleFunc("/", s.handleProxy)
	fmt.Printf("Starting HTTP proxy server on %s...\n", proxyAddr)
	err := http.ListenAndServe(proxyAddr, nil)
	if err != nil {
		fmt.Printf("Error starting HTTP proxy server: %s\n", err)
	}
}
