package main

import (
	"encoding/json"
	"karinto/trx-downloader/cache"
	"karinto/trx-downloader/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMonzoRefreshToken(t *testing.T) {

	config.InitConfig()

	cache.WriteCache(cache.MonzoAccessTokenKey, "accessToken_before")
	cache.WriteCache(cache.MonzoRefreshTokenKey, "refreshToken_before")

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			decoder := json.NewDecoder(r.Body)
			var payload map[string]string
			err := decoder.Decode(&payload)
			if payload["refresh_token"] != "refreshToken_before" {
				t.Errorf("Failed to read request body: %v", err)
			}

			w.WriteHeader(http.StatusOK)
			_, err = w.Write([]byte(`{
                "access_token": "accessToken_after",
                "refresh_token": "refreshToken_after"
            }`))
			if err != nil {
				t.Errorf("Failed to write json response: %v", err)
			}
		}))

	config.Config["monzo_refresh_url"] = server.URL

	RefreshToken()

	if cache.ReadCache(cache.MonzoAccessTokenKey) != "accessToken_after" && cache.ReadCache(cache.MonzoRefreshTokenKey) != "refreshToken_after" {
		t.Error("Failed to refresh token")
	}

}
