package cache

import (
	"os"
	"reflect"
	"testing"
)

func TestCacheReadWrite(t *testing.T) {
    prepareCacheFile(t)

	key := MonzoAccessTokenKey
	value := "xxxxx"

	isWritten := writeCache(key, value)

	if !isWritten {
		t.Error("Failed to write cache")
	}

	result := readCache(key)

	if result != value {
		t.Errorf("Read value %s is different from written value %s", result, value)
	}
}

func TestCacheFileWriteRead(t *testing.T) {
    prepareCacheFile(t)
	data := map[string]string{"testKey": "testValue"}

	writeOnFile(data)

	result := readFromFile()

	if !reflect.DeepEqual(data, result) {
		t.Error("Failed file test")
	}

}

func TestCacheReadNilValue(t *testing.T) {
    prepareCacheFile(t)

	// retrive a key that does not exist
	result := readCache("KeyWithoutValue")

	if result != "" {
		t.Error("Failed config read nil test")
	}

}

func prepareCacheFile(t *testing.T) {
	f, err := os.CreateTemp("", "test_config.json")
    if err != nil {
        t.Errorf("Error on creating tmp cache: %v", err)
    }

    defer os.Remove(f.Name())

    config["cache_file_path"] = f.Name()
}
