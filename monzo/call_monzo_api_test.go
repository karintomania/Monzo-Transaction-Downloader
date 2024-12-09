package monzo

import (
	"encoding/json"
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

func mockServer(t *testing.T, accessToken string, accountId string, tr TransactionResponse) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// assert URL
			if got, want := r.URL.Path, "/transactions"; got != want {
				t.Errorf("Expected to request %s, got: %s", want, got)
			}

			if got, want := r.URL.Query().Get("expand[]"), "merchant"; got != want {
				t.Errorf("Expected expand[] to be %s, got: %s", want, got)
			}

			if got, want := r.URL.Query().Get("account_id"), accountId; got != want {
				t.Errorf("Expected expand[] to be %s, got: %s", want, got)
			}

			// assert Header
			if got := r.Header.Get("Authorization"); got != accessToken {
				t.Errorf("Expected %s, but got %s", accessToken, got)
			}

			jsonResponse, err := json.Marshal(tr)
			if err != nil {
				t.Errorf("Failed to marshal response: %v", err)
			}

			_, err = w.Write(jsonResponse)
			if err != nil {
				t.Errorf("Failed to write json response: %v", err)
			}
		}))

	return server
}

func TestMonzoDownloadTransaction(t *testing.T) {

	accessToken := "accessToken_123"
	cache.Write(cache.MonzoAccessTokenKey, accessToken)

	accountId := "accountId_123"
	config.Set(config.MONZO_ACCOUNT_ID, accountId)

	wantResponse := TransactionResponse{
		[]Transaction{
			{
				ID:          "tx_000001",
				Amount:      -1000,
				Description: "description 1",
				Merchant: Merchant{
					Name: "merchant A",
				},
			},
		},
	}

	server := mockServer(t, accessToken, accountId, wantResponse)

	config.Set(
		config.MONZO_TRANSACTIONS_URL,
		fmt.Sprintf("%s/transactions?expand[]=merchant&account_id=%s", server.URL, accountId),
	)

	result := DownloadTransactions()

	for i, want := range wantResponse.Transactions {
		if want.ID != result[i].ID {
			t.Errorf("Expected %s, but got %s", result[i].ID, want.ID)
		}
		if want.Amount != result[i].Amount {
			t.Errorf("Expected %d, but got %d", result[i].Amount, want.Amount)
		}
		if want.Description != result[i].Description {
			t.Errorf("Expected %s, but got %s", result[i].Description, want.Description)
		}
		if want.Merchant.Name != result[i].Merchant.Name {
			t.Errorf("Expected %s, but got %s", result[i].Merchant.Name, want.Merchant.Name)
		}
	}

}
