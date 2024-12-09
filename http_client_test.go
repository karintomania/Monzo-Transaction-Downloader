package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpGet(t *testing.T) {

	header := map[string]string{"Authorization": "Bearer 12345"}

	expectedResponseBody := []byte(`{"data":"success"}`)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-path" {
			t.Errorf("Expected to request '/test-path', got: %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected to GET method, got: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer 12345" {
			t.Errorf("Expected 'Authorization: Bearer 12345', got: %s", r.Header.Get("Authorization"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write(expectedResponseBody)
	}))

	url := server.URL + "/test-path"

	defer server.Close()

	actual, _ := httpGet(url, header)

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
		w.Write(expectedResponseBody)
	}))

	url := server.URL + "/test-path"

	defer server.Close()

	actual, _ := httpPost(url, header, payload)

	if string(expectedResponseBody) != string(actual) {
		t.Errorf("Expected %s, got %s", expectedResponseBody, actual)
	}

}