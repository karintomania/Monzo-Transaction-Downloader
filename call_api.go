package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// Define structs to match the JSON structure
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	// ExpiresIn string `"json:"expires_in"`
	// ClientId string `"json:"client_id"`
	// Scope string `"json:"scope"`
	// TokenType string `"json:"token_type"`
	// UserId string `"json:"user_id"`
}

func callHttpRefreshToken() []byte {
	// Make the HTTP GET request
	req, err := http.NewRequest("GET", "https://api.monzo.com/oauth2/token", nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+readConfig(MonzoRefreshTokenKey))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read request body: %v", err)
	}

	return b
}

func callHttpRefreshTokenMock() []byte {

	return []byte(`{
    "access_token":"accessToken_xxx",
    "refresh_token":"refreshToken_xxx",
    "expires_in":21600,
    "client_id":"client_123",
    "scope":"accounts",
    "token_type":"refresh",
    "user_id":"user_123"
    }`)
}

func callRefreshToken() RefreshTokenResponse {

	body := callHttpRefreshTokenMock()

	log.Println(string(body))

	// Parse the JSON response
	var refreshTokenResponse RefreshTokenResponse
	if err := json.Unmarshal(body, &refreshTokenResponse); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	log.Println(refreshTokenResponse)

	return refreshTokenResponse

}
