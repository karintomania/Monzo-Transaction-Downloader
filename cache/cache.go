package cache

import (
	"encoding/json"
	"karinto/trx-downloader/config"
	"log"
	"os"
)

const MonzoAccessTokenKey = "MonzoAccessToken"
const MonzoRefreshTokenKey = "MonzoRefreshToken"
const MonzoClientIdKey = "MonzoClientId"
const MonzoClientSecretKey = "MonzoClientSecret"

func Write(key, value string) bool {

	c := readFromFile()
	c[key] = value

	writeOnFile(c)

	return true
}

func Read(key string) string {
	c := readFromFile()

	return c[key]
}

func writeOnFile(data map[string]string) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal json: %v", err)
	}

	path := config.Get("cache_file_path")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatalf("Failed to open cache file: %v", err)
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		log.Fatalf("Failed to write cache file: %v", err)
	}
}

func readFromFile() map[string]string {
	path := config.Get("cache_file_path")

	err := ensureFileExists(path)
	if err != nil {
		log.Fatalf("File doesn't exists on %s: %v", path, err)
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Can't open file %s: %v", path, err)
	}

	decoder := json.NewDecoder(f)

	cache := make(map[string]string)

	err = decoder.Decode(&cache)
	if err != nil {
		log.Fatalf("Can't decode cache file %s: %v", path, err)
	}

	return cache
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
