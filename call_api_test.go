package main

import (
	"testing"
)

func TestCall(t *testing.T) {
	t.Log("aaaa")
	users := call()
	user := users[0]
	t.Log(user.ID)
}
