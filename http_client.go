package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func httpGet(url string, header map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return sendRequest(header, req)
}

func httpPostForm(urlString string, header map[string]string, formDataMap map[string][]string) ([]byte, error) {

    formData := url.Values(formDataMap)
    req, err := http.NewRequest("POST", urlString, strings.NewReader(formData.Encode()))
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

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != 200 {
        return nil, fmt.Errorf("unexpected status code: %d, %s", resp.StatusCode, string(body[:]))
    }

	return body, nil
}
