package monzo

import (
	"fmt"
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
			if r.URL.Path != "/oauth2/token" {
				t.Errorf("Expected to request '/oauth2/token', got: %s", r.URL.Path)
			}

			if err := r.ParseForm(); err != nil {
				t.Errorf("Failed to parse form data: %v", err)
			}

			if r.FormValue("refresh_token") != "refreshToken_before" {
				t.Errorf("Expected refreshToken_before, but got %s", r.FormValue("refresh_token"))
			}

			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{
                "access_token": "accessToken_after",
                "refresh_token": "refreshToken_after"
            }`))
			if err != nil {
				t.Errorf("Failed to write json response: %v", err)
			}
		}))

	config.Set("monzo_refresh_url", fmt.Sprintf("%s/oauth2/token", server.URL))

	RefreshToken()

	if cache.Read(cache.MonzoAccessTokenKey) != "accessToken_after" &&
		cache.Read(cache.MonzoRefreshTokenKey) != "refreshToken_after" {
		t.Error("Failed to refresh token")
	}

}

func TestMonzoDownloadTransaction(t *testing.T) {

	accessToken := "accessToken_123"
	cache.Write(cache.MonzoAccessTokenKey, accessToken)

	accountId := "accountId_123"
	config.Set(config.MONZO_ACCOUNT_ID, accountId)

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// assert URL
			wantUrl := fmt.Sprintf("/transactions?expand[]=merchant&account_id=%s", accountId)
			if got := r.URL.Path; got != wantUrl {
				t.Errorf("Expected to request %s, got: %s", wantUrl, got)
			}

			// assert Header
			if got := r.Header.Get("Authorization"); got != "accessToken" {
				t.Errorf("Expected %s, but got %s", accessToken, got)
			}

			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{
  "transactions": [
        {
           "id": "tx_0000AnfhaSVKiaF0Y9x6vZ",
           "created": "2024-11-03T16:47:18.081Z",
           "description": "Description xxx",
           "amount": -1640,
           "currency": "GBP",
           "merchant": {
             "name": "Merchant A",
           },
           "local_amount": -1640,
           "local_currency": "GBP",
           "updated": "2024-11-04T08:45:33.734Z",
        }
    ]
}`))
			if err != nil {
				t.Errorf("Failed to write json response: %v", err)
			}
		}))

	config.Set(config.MONZO_TRANSACTIONS_URL, server.URL)

	result := DownloadTransactions()

	transaction := result[0]

	if transaction.id != "tx_0000AnfhaSVKiaF0Y9x6vZ" {
		t.Errorf("Expected tx_0000AnfhaSVKiaF0Y9x6vZ, but got %s", transaction.id)
	}
	if transaction.amount != -1640 {
		t.Errorf("Expected -1640, but got %d", transaction.amount)
	}
	if transaction.merchantName != "Merchant A" {
		t.Errorf("Expected Merchant A, but got %s", transaction.merchantName)
	}
}
