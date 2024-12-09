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

	b, err := json.Marshal(data)
	log.Println(string(b))

	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Dir(os.Args[0]) + "/config.json"
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0777)
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
	path := filepath.Dir(os.Args[0]) + "/config.json"

	err := ensureFileExists(path)

	if err != nil {
		log.Fatal(err)
	}

	body, err := os.ReadFile(path)

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
