package monzo

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"karinto/trx-downloader/cache"
	"karinto/trx-downloader/config"
	"karinto/trx-downloader/httpClient"
	"log"
	"os"
	"time"
)

const DateLayout = "2006-01-02T15:04:05.000Z"

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
	ID          string   `json:"id"`
	Created     string   `json:"created"`
	Description string   `json:"description"`
	Amount      int      `json:"amount"`
	Currency    string   `json:"currency"`
	Merchant    Merchant `json:"merchant"`
	Notes       string   `json:"notes"`
}

type Merchant struct {
	Name string `json:"name"`
}

func DownloadTransactions(since time.Time) []Transaction {

	baseUrl := config.Get(config.MONZO_TRANSACTIONS_URL)

	accessToken := cache.Read(cache.MonzoAccessTokenKey)

	header := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Accept":        "application/json",
		"Authorization": "Bearer " + accessToken,
	}

	accountId := config.Get(config.MONZO_ACCOUNT_ID)

	queries := map[string]string{
		"expand[]":   "merchant",
		"account_id": accountId,
		"since":      since.Format(DateLayout),
	}

	body, err := httpClient.Get(baseUrl, header, queries)
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

// if successfule, return path for the CSV
func EncodeTransactionsCsv(transactions []Transaction) (string, error) {

	// Get current working directory
	path, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Failed to get current working directory: %w", err)
	}

	path = path + "result_1.csv"

	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("Failed to create file: %w", err)
	}

	defer f.Close()

	writer := csv.NewWriter(f)

	defer writer.Flush()

	header := []string{
		"ID",
		"Created",
		"Description",
		"Amount",
		"Currency",
		"Merchant",
		"Notes",
	}

	var rows [][]string

	for _, transaction := range transactions {
		row := []string{
			transaction.ID,
			transaction.Created,
			transaction.Description,
			fmt.Sprintf("%d", transaction.Amount),
			transaction.Currency,
			transaction.Merchant.Name,
			transaction.Notes,
		}

		rows = append(rows, row)
	}

	err = writer.Write(header)
	if err != nil {
		return "", fmt.Errorf("Failed to write header: %w", err)
	}

	for _, row := range rows {
		if err = writer.Write(row); err != nil {
			return "", fmt.Errorf("Failed to write row: %w", err)
		}
	}

	return path, nil
}
