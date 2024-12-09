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
	cache.Write(cache.MonzoAccessTokenKey, refreshTokenResponse.AccessToken)
	cache.Write(cache.MonzoRefreshTokenKey, refreshTokenResponse.RefreshToken)
}

func callRefreshToken() RefreshTokenResponse {
	// Make the HTTP GET request
	url := config.Config["monzo_refresh_url"]

	refreshToken := cache.Read(cache.MonzoRefreshTokenKey)

	header := map[string]string{"Content-Type": "application/x-www-form-urlencoded", "Accept": "application/json"}

	formDataMap := map[string][]string{
		"grant_type":    {"refresh_token"},
		"client_id":     {config.Config["monzo_client_id"]},
		"client_secret": {config.Config["monzo_client_secret"]},
		"refresh_token": {refreshToken},
	}

	body, err := httpPostForm(url, header, formDataMap)
	if err != nil {
		log.Fatalf("Failed to make HTTP request: %v", err)
	}

	// Parse the JSON response
	var refreshTokenResponse RefreshTokenResponse
	if err := json.Unmarshal(body, &refreshTokenResponse); err != nil {
		log.Fatalf("Failed to parse JSON %s: %v", string(body), err)
	}

	return refreshTokenResponse

}

func GetTransactions() {
}
