package config

import (
	"encoding/json"
	"os"
	"testing"
)

func TestReadConfigFile(t *testing.T) {

	configWant := map[string]string{"testKey1": "testValue1", "testKey2": "testValue2"}
	bytes, err := json.Marshal(configWant)
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

	for key, want := range configWant {
		if value, exists := Config[key]; !exists || value != want {
			t.Errorf("Config mismatch for key %s: expected %s, got %s", key, want, value)
		}
	}

}
