package httpClient

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHttpGet(t *testing.T) {

	header := map[string]string{"Authorization": "Bearer 12345"}

	queries := map[string]string{"before": "2024-01-01T00:00:00Z"}

	expectedResponseBody := []byte(`{"data":"success"}`)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-path" {
			t.Errorf("Expected to request '/test-path', got: %s", r.URL.Path)
		}

		if r.Method != "GET" {
			t.Errorf("Expected to GET method, got: %s", r.URL.Path)
		}

		gotQueries := r.URL.Query()
		for k, want := range queries {
			got := gotQueries.Get(k)
			if got != want {
				t.Errorf("Expected query value %s for %s, got: %s", want, k, got)
			}
		}

		if r.Header.Get("Authorization") != "Bearer 12345" {
			t.Errorf("Expected 'Authorization: Bearer 12345', got: %s", r.Header.Get("Authorization"))
		}

		w.WriteHeader(http.StatusOK)
		_, err := w.Write(expectedResponseBody)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
	}))

	url := server.URL + "/test-path"

	defer server.Close()

	actual, _ := Get(url, header, queries)

	if string(expectedResponseBody) != string(actual) {
		t.Errorf("Expected %s, got %s", expectedResponseBody, actual)
	}

}

func TestHttpPost(t *testing.T) {

	header := map[string]string{"Authorization": "Bearer 12345"}
	expectedResponseBody := []byte(`{"data":"success"}`)
	payload := []byte(`{"key": "value"}`)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		payloadReceived, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error on reading payload: %v", err)
		}
		if string(payloadReceived) != `{"key": "value"}` {
			t.Errorf("Expected to request '/test-path', got: %s", r.URL.Path)
		}

		if r.URL.Path != "/test-path" {
			t.Errorf("Expected to request '/test-path', got: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected to POST method, got: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer 12345" {
			t.Errorf("Expected 'Authorization: Bearer 12345', got: %s", r.Header.Get("Authorization"))
		}

		w.WriteHeader(http.StatusOK)

		_, err = w.Write(expectedResponseBody)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
	}))

	url := server.URL + "/test-path"

	defer server.Close()

	actual, _ := Post(url, header, payload)

	if string(expectedResponseBody) != string(actual) {
		t.Errorf("Expected %s, got %s", expectedResponseBody, actual)
	}

}

func TestHttpPostForm(t *testing.T) {

	header := map[string]string{"Authorization": "Bearer 12345"}
	expectedResponseBody := []byte(`{"data":"success"}`)
	formDataMap := map[string][]string{"key1": {"value1"}, "key2": {"value2"}}
	formData := url.Values(formDataMap)
	formDataEncoded := formData.Encode()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		payloadReceived, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error on reading payload: %v", err)
		}

		if string(payloadReceived) != formDataEncoded {
			t.Errorf("Expected formatData: %s, got: %s", formDataEncoded, payloadReceived)
		}

		if r.URL.Path != "/test-path" {
			t.Errorf("Expected to request '/test-path', got: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected to POST method, got: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer 12345" {
			t.Errorf("Expected 'Authorization: Bearer 12345', got: %s", r.Header.Get("Authorization"))
		}

		w.WriteHeader(http.StatusOK)

		_, err = w.Write(expectedResponseBody)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
	}))

	url := server.URL + "/test-path"

	defer server.Close()

	actual, _ := PostForm(url, header, formData)

	if string(expectedResponseBody) != string(actual) {
		t.Errorf("Expected %s, got %s", expectedResponseBody, actual)
	}

}
