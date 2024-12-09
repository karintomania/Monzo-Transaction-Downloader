package main

import (
	"encoding/json"
	"karinto/trx-downloader/cache"
	"karinto/trx-downloader/config"
	"log"
)

// Define structs to match the JSON structure
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func RefreshToken() {
    // Call the function to refresh the token
    refreshTokenResponse := callRefreshToken()
    // Write the new access token to the cache
    cache.WriteCache(cache.MonzoAccessTokenKey, refreshTokenResponse.AccessToken)
    cache.WriteCache(cache.MonzoRefreshTokenKey, refreshTokenResponse.RefreshToken)
}

func callRefreshToken() RefreshTokenResponse {
	// Make the HTTP GET request
	url := config.Config["monzo_refresh_url"]

	refreshToken := cache.ReadCache(cache.MonzoRefreshTokenKey)

	header := map[string]string{"Authorization": refreshToken}

	body, err := httpPost(url, header, nil)
	if err != nil {
		log.Fatalf("Failed to make HTTP request: %v", err)
	}

	// Parse the JSON response
	var refreshTokenResponse RefreshTokenResponse
	if err := json.Unmarshal(body, &refreshTokenResponse); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	return refreshTokenResponse

}
