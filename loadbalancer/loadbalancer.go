package loadbalancer

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
)

type LoadBalancer struct {
	Nodes []*Node
	mu    sync.Mutex
}

func NewLoadBalancer(nodes []*Node) *LoadBalancer {
	return &LoadBalancer{
		Nodes: nodes,
	}
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read the request body to determine its size
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()
	bodySize := len(body)

	// Find a suitable node
	lb.mu.Lock()
	defer lb.mu.Unlock()
	for _, node := range lb.Nodes {
		if node.AllowRequest(bodySize) {
			// Forward the request to the selected node
			req, err := http.NewRequest(r.Method, node.Address, io.NopCloser(io.LimitReader(io.MultiReader(io.NopCloser(io.Reader(nil)), nil), int64(len(body)))))
			if err != nil {
				http.Error(w, "Bad Gateway", http.StatusBadGateway)
				return
			}

			// Copy headers
			for key, values := range r.Header {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}

			// Set the body
			req.Body = io.NopCloser(io.Reader(io.MultiReader(io.NopCloser(io.Reader(nil)), nil)))

			// Perform the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				http.Error(w, "Bad Gateway", http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()

			// Copy the headers
			for key, values := range resp.Header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}

			// Write the status code
			w.WriteHeader(resp.StatusCode)

			// Copy the body to the ResponseWriter
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				log.Printf("Error copying response body: %v", err)
			}
			return
		}
	}

	// If no nodes are available
	http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
}

func forwardRequest(address string, originalReq *http.Request, w http.ResponseWriter, body []byte) error {
	// Create a new request to the backend
	parsedURL, err := url.Parse(address)
	if err != nil {
		return err
	}
	proxyReq, err := http.NewRequest(originalReq.Method, parsedURL.String()+originalReq.RequestURI, bytes.NewReader(body))
	if err != nil {
		return err
	}

	// Copy headers
	for key, values := range originalReq.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Write the status code
	w.WriteHeader(resp.StatusCode)

	// Copy the response body to the ResponseWriter
	_, err = io.Copy(w, resp.Body)
	return err
}
