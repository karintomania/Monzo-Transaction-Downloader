package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func httpGet(url string, header map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return sendRequest(header, req)
}

func httpPost(url string, header map[string]string, payload []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return sendRequest(header, req)
}

func sendRequest(header map[string]string, req *http.Request) ([]byte, error) {
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
