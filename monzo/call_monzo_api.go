package monzo

import (
	"encoding/json"
	"karinto/trx-downloader/cache"
	"karinto/trx-downloader/config"
	"karinto/trx-downloader/httpClient"
	"log"
	"time"
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
	url := config.Get("monzo_refresh_url")

	refreshToken := cache.Read(cache.MonzoRefreshTokenKey)

	header := map[string]string{"Content-Type": "application/x-www-form-urlencoded", "Accept": "application/json"}

	formDataMap := map[string][]string{
		"grant_type":    {"refresh_token"},
		"client_id":     {config.Get("monzo_client_id")},
		"client_secret": {config.Get("monzo_client_secret")},
		"refresh_token": {refreshToken},
	}

	body, err := httpClient.PostForm(url, header, formDataMap)
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

type TransactionResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	ID                  string         `json:"id"`
	Created             time.Time      `json:"created"`
	Description         string         `json:"description"`
	Amount              int            `json:"amount"`
	Currency            string         `json:"currency"`
	Merchant            Merchant       `json:"merchant"`
	Category            string         `json:"category"`
	Categories          map[string]int `json:"categories"`
	Settled             time.Time      `json:"settled"`
	LocalAmount         int            `json:"local_amount"`
	LocalCurrency       string         `json:"local_currency"`
	Updated             time.Time      `json:"updated"`
}

type Merchant struct {
	ID       string `json:"id"`
	GroupID  string `json:"group_id"`
	Name     string `json:"name"`
	Logo     string `json:"logo"`
	Category string `json:"category"`
}

func DownloadTransactions() []Transaction {

	url := config.Get(config.MONZO_TRANSACTIONS_URL)

	accessToken := cache.Read(cache.MonzoAccessTokenKey)

	header := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Accept":        "application/json",
		"Authorization": accessToken,
	}

	body, err := httpClient.Get(url, header)
	if err != nil {
		log.Fatalf("Failed to make HTTP request for monzo transaction download: %v", err)
	}

	// var response transactionResponse
	var tr TransactionResponse
	if err = json.Unmarshal(body, &tr); err != nil {
		log.Fatalf("Failed to unmarshal response body: %v", err)
	}

	return tr.Transactions
}
