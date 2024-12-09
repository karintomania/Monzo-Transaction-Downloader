package cache

import (
	"os"
	"reflect"
	"testing"
)

func TestCacheReadWrite(t *testing.T) {
	path := Fake()
	defer os.Remove(path)

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
	path := Fake()
	defer os.Remove(path)

	data := map[string]string{"testKey": "testValue"}

	writeOnFile(data)

	result := readFromFile()

	if !reflect.DeepEqual(data, result) {
		t.Error("Failed file test")
	}

}

func TestCacheReadNilValue(t *testing.T) {
	path := Fake()
	defer os.Remove(path)

	// retrive a key that does not exist
	result := Read("KeyWithoutValue")

	if result != "" {
		t.Error("Failed config read nil test")
	}

}
