package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var config map[string]string

func Get(key string) string {
	if config == nil {
		initConfig()
	}

	return config[key]
}

// only use for testing
func Set(key string, value string) {
	if config == nil {
		initConfig()
	}

	config[key] = value
}

func initConfig() {
	config = make(map[string]string)
	configDir := getConfigFolderPath()
	config["config_file_path"] = filepath.Join(configDir, "/config.json")
	config["cache_file_path"] = filepath.Join(configDir, "/cache.json")

	readConfigFile(config["config_file_path"])
}

func readConfigFile(path string) {
	f, err := os.Open(path)

	if err != nil {
		log.Fatalf("Failed to read %s. The error happend: %v", path, err)
	}

	decoder := json.NewDecoder(f)

	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to decode config: %v", err)
	}
}

func getConfigFolderPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".config", "trx-downloader")

	err = os.MkdirAll(configDir, 0775)
	if err != nil {
		log.Fatalf("Error creating config directory: %v", err)
	}

	return configDir
}
