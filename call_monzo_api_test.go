package main

import (
	"karinto/trx-downloader/cache"
	"karinto/trx-downloader/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMonzoRefreshToken(t *testing.T) {

	cache.Write(cache.MonzoAccessTokenKey, "accessToken_before")
	cache.Write(cache.MonzoRefreshTokenKey, "refreshToken_before")

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			err := r.ParseForm()
			if err != nil {
				t.Errorf("Failed to parse form data: %v", err)
			}

			if r.FormValue("refresh_token") != "refreshToken_before" {
				t.Errorf("Expected refreshToken_before, but got %s", r.FormValue("refresh_token"))
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

	config.Set("monzo_refresh_url", server.URL)

	RefreshToken()

	if cache.Read(cache.MonzoAccessTokenKey) != "accessToken_after" &&
		cache.Read(cache.MonzoRefreshTokenKey) != "refreshToken_after" {
		t.Error("Failed to refresh token")
	}

}
