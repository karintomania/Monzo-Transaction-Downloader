package httpClient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func Get(baseUrl string, header map[string]string, queries map[string]string) ([]byte, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	query := u.Query()

	for k, v := range queries {
		query.Set(k, v)
	}

	u.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return sendRequest(header, req)
}

func PostForm(urlString string, header map[string]string, formDataMap map[string][]string) ([]byte, error) {

	formData := url.Values(formDataMap)
	req, err := http.NewRequest("POST", urlString, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return sendRequest(header, req)
}

func Post(url string, header map[string]string, payload []byte) ([]byte, error) {
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
