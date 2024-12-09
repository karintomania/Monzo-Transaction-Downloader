package config

import (
	"encoding/json"
	"os"
	"testing"
)

func mockConfigFile(wantConfig *map[string]string, t *testing.T) {
	bytes, err := json.Marshal(wantConfig)
	if err != nil {
		t.Fatalf("Failed to marshal configExpected: %v", err)
	}

	f, err := os.CreateTemp("", "test_config.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer os.Remove(f.Name())

	_, err = f.Write(bytes)
	if err != nil {
		t.Fatalf("Failed to write on config file: %v", err)
	}

	readConfigFile(f.Name())
}

func TestReadConfigFile(t *testing.T) {

	wantConfig := map[string]string{"testKey1": "testValue1", "testKey2": "testValue2"}
	mockConfigFile(&wantConfig, t)

	for key, want := range wantConfig {
		if got := Get(key); got != want {
			t.Errorf("Config mismatch for key %s: expected %s, got %s", key, want, got)
		}
	}

}
