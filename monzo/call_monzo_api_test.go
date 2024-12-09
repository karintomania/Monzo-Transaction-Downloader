package monzo

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"karinto/trx-downloader/cache"
	"karinto/trx-downloader/config"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
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

func mockServer(t *testing.T, accessToken string, tr TransactionResponse, q map[string]string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// assert URL
			if got, want := r.URL.Path, "/transactions"; got != want {
				t.Errorf("Expected to request %s, got: %s", want, got)
			}

			gotQs := r.URL.Query()
			for k, want := range q {
				if got := gotQs.Get(k); got != want {
					t.Errorf("Expected query %s to be %s, got: %s", k, want, got)
				}
			}

			// assert Auth
			if got, want := r.Header.Get("Authorization"), "Bearer "+accessToken; got != want {
				t.Errorf("Expected Authorization %s, but got %s", want, got)
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
				Created:     "2024-01-02T03:04:05.006Z",
				Amount:      -1000,
				Currency:    "GBP",
				Notes:       "test note 1",
				Description: "description 1",
				Merchant: Merchant{
					Name: "merchant 1",
				},
			},
			{
				ID:          "tx_000002",
				Created:     "2025-02-03T04:05:06.007Z",
				Amount:      -1001,
				Currency:    "USD",
				Notes:       "test note 2",
				Description: "description 2",
				Merchant: Merchant{
					Name: "merchant 2",
				},
			},
		},
	}

	since := time.Now().AddDate(0, -1, 0)

	wantQueries := map[string]string{
		"expand[]":   "merchant",
		"account_id": accountId,
		"since":      since.Format(DateLayout),
	}

	server := mockServer(t, accessToken, wantResponse, wantQueries)

	config.Set(
		config.MONZO_TRANSACTIONS_URL,
		fmt.Sprintf("%s/transactions?expand[]=merchant&account_id=%s", server.URL, accountId),
	)

	result := DownloadTransactions(since)

	for i, want := range wantResponse.Transactions {
		if want.ID != result[i].ID {
			t.Errorf("Expected %s, but got %s", result[i].ID, want.ID)
		}
		if want.Created != result[i].Created {
			t.Errorf("Expected %s, but got %s", result[i].Created, want.Created)
		}
		if want.Amount != result[i].Amount {
			t.Errorf("Expected %d, but got %d", result[i].Amount, want.Amount)
		}
		if want.Currency != result[i].Currency {
			t.Errorf("Expected %s, but got %s", result[i].Currency, want.Currency)
		}
		if want.Notes != result[i].Notes {
			t.Errorf("Expected %s, but got %s", result[i].Notes, want.Notes)
		}
		if want.Description != result[i].Description {
			t.Errorf("Expected %s, but got %s", result[i].Description, want.Description)
		}
		if want.Merchant.Name != result[i].Merchant.Name {
			t.Errorf("Expected %s, but got %s", result[i].Merchant.Name, want.Merchant.Name)
		}
	}
}

func TestEncodeTransactionsCsv(t *testing.T) {
	transactions := []Transaction{
		{
			ID:          "tx_000001",
			Created:     "2024-01-02T03:04:05.006Z",
			Amount:      -1000,
			Currency:    "GBP",
			Notes:       "test note 1",
			Description: "description 1",
			Merchant: Merchant{
				Name: "merchant 1",
			},
		},
		{
			ID:          "tx_000002",
			Created:     "2025-02-03T04:05:06.007Z",
			Amount:      -1001,
			Currency:    "USD",
			Notes:       "test note 2",
			Description: "description 2",
			Merchant: Merchant{
				Name: "merchant 2",
			},
		},
	}

	path, err := EncodeTransactionsCsv(transactions)
	if err != nil {
		t.Errorf("Failed to encode transactions to CSV: %v", err)
	}
    defer os.Remove(path)

	f, err := os.Open(path)
	if err != nil {
		t.Errorf("Failed to open CSV: %v", err)
	}
	reader := csv.NewReader(f)

	rowsGot, err := reader.ReadAll()
	if err != nil {
		t.Errorf("Failed to open CSV: %v", err)
	}

	for i, got := range rowsGot {
		// Assert header
		if i == 0 {
			if got[0] != "ID" ||
				got[1] != "Created" ||
				got[2] != "Description" ||
				got[3] != "Amount" ||
				got[4] != "Currency" ||
				got[5] != "Merchant" ||
				got[6] != "Notes" {
				t.Errorf("Expected header, but got %v", got)
			}
			continue
		}

		// assert rows
		want := transactions[i-1]

		if got[0] != want.ID ||
			got[1] != want.Created ||
			got[2] != want.Description ||
			got[3] != fmt.Sprintf("%d", want.Amount) ||
			got[4] != want.Currency ||
			got[5] != want.Merchant.Name ||
			got[6] != want.Notes {
			t.Errorf("Expected %v, but got %v", want, got)
		}
	}



}
