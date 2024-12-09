package cache

import (
	"karinto/trx-downloader/config"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCacheReadWrite(t *testing.T) {
	prepareCacheFile(t)

	key := MonzoAccessTokenKey
	value := "xxxxx"

	t.Log("test")
	isWritten := Write(key, value)

	if !isWritten {
		t.Error("Failed to write cache")
	}

	t.Log("test")
	result := Read(key)

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
	result := Read("KeyWithoutValue")

	if result != "" {
		t.Error("Failed config read nil test")
	}

}

func prepareCacheFile(t *testing.T) {
	cleanupCacheFile(t)

	f, err := os.CreateTemp("", "test_config.json")
	if err != nil {
		t.Errorf("Error on creating tmp cache: %v", err)
	}

	_, err = f.Write([]byte("{}"))
	if err != nil {
		t.Errorf("Error on writing tmp cache: %v", err)
	}

	config.Set("cache_file_path", f.Name())
}

func cleanupCacheFile(t *testing.T) {
	files, err := filepath.Glob("/tmp/test_config.json*")
	if err != nil {
		t.Errorf("Error on deleting tmp cache: %v", err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			t.Errorf("Error on deleting tmp cache: %v", err)
		}
	}
}
