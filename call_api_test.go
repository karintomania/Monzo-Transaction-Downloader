package main

import (
	"testing"
)

func TestCallRefreshToken(t *testing.T) {
	result := callRefreshToken()

	if result.AccessToken != "accessToken_xxx" && result.RefreshToken != "refreshToken_xxx" {
		t.Error("Failed to refresh token")
	}
}
