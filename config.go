package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const MonzoAccessTokenKey = "MonzoAccessToken"
const MonzoRefreshTokenKey = "MonzoRefreshToken"
const MonzoClientIdKey = "MonzoClientId"
const MonzoClientSecretKey = "MonzoClientSecret"

func writeConfig(key, value string) bool {

	c := readFromFile()
	c[key] = value

	writeOnFile(c)

	return true
}

func readConfig(key string) string {
	c := readFromFile()

	return c[key]
}

func writeOnFile(data map[string]string) {
	log.SetFlags(log.Lshortfile)
	b, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	path := getConfigFilePath()
    f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Write(b)

	if err != nil {
		log.Fatal(err)
	}
}

func readFromFile() map[string]string {
	path := getConfigFilePath()

	err := ensureFileExists(path)

	if err != nil {
		log.Fatal(err)
	}

	body, err := os.ReadFile(path)

	log.Println(string(body))

	if err != nil {
		log.Fatal(err)
	}

	if len(body) == 0 {
		body = []byte("{}")
	}

	var config map[string]string
	err = json.Unmarshal(body, &config)

	if err != nil {
		log.Fatal(err)
	}

	return config
}

func ensureFileExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

func getConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v", err)
	}

    configDir := filepath.Join(homeDir, ".config", "trx-downloader")

    err = os.MkdirAll(configDir, 0775)
	if err != nil {
		log.Fatalf("Error creating config directory: %v", err)
	}
    configFilePath := filepath.Join(configDir, "/config.json")

    return configFilePath
}
