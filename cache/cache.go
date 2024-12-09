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

func WriteCache(key, value string) bool {
	c := readFromFile()
	c[key] = value

	writeOnFile(c)

	return true
}

func ReadCache(key string) string {
	c := readFromFile()

	return c[key]
}

func writeOnFile(data map[string]string) {
	b, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	path := config.Config["cache_file_path"]
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
	path := config.Config["cache_file_path"]

	err := ensureFileExists(path)

	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

    decoder := json.NewDecoder(f)

    cache := make(map[string]string)

	decoder.Decode(&cache)

	if err != nil {
		log.Fatal(err)
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
