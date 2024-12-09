package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type geometry interface {
	area() float64
	perim() float64
}
type httpClient interface {
	httpGet(url string, header map[string]string) ([]byte, error)
	httpPost(url string, header map[string]string, payload []byte) ([]byte, error)
}

type httpClientImpl struct{}

func (h httpClientImpl) httpGet(url string, header map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return h.sendRequest(header, req)
}

func (h httpClientImpl) httpPost(url string, header map[string]string, payload []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return h.sendRequest(header, req)
}

func (h httpClientImpl) sendRequest(header map[string]string, req *http.Request) ([]byte, error) {
	// Add headers to the request
	for k, v := range header {
		req.Header.Add(k, v)
	}

	// Initialize an HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
