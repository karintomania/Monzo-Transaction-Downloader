package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var Config map[string]string = make(map[string]string)

func InitConfig() {
	configDir := getConfigFolderPath()
	Config["config_file_path"] = filepath.Join(configDir, "/config.json")
	Config["cache_file_path"] = filepath.Join(configDir, "/cache.json")

	readConfigFile(Config["config_file_path"])
}

func readConfigFile(path string) {
	f, err := os.Open(path)

	if err != nil {
		log.Fatalf("Failed to read %s. The error happend: %v", path, err)
	}

	decoder := json.NewDecoder(f)

	decoder.Decode(&Config)
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
