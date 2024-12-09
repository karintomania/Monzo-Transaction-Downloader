package main

import (
	"reflect"
	"testing"
)

func TestConfigReadWrite(t *testing.T) {
	key := MonzoAccessTokenKey
	value := "xxxxx"

	isWritten := writeConfig(key, value)

	if !isWritten {
		t.Error("Failed to write cache")
	}

	result := readConfig(key)

	if result != value {
		t.Errorf("Read value %s is different from written value %s", result, value)
	}
}

func TestConfigFileWriteRead(t *testing.T) {
	data := map[string]string{"testKey": "testValue"}

	writeOnFile(data)

	result := readFromFile()

	if !reflect.DeepEqual(data, result) {
		t.Error("Failed file test")
	}

}
