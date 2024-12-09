package main

import (
	"karinto/trx-downloader/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCallRefreshToken(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
            "access_token": "accessToken_xxx",
            "refresh_token": "refreshToken_xxx"
            }`))
		}))

    config.InitConfig()
	config.Config["monzo_refresh_url"] = server.URL

	result := callRefreshToken()

	if result.AccessToken != "accessToken_xxx" && result.RefreshToken != "refreshToken_xxx" {
		t.Error("Failed to refresh token")
	}
}
