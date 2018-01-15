package main

import (
	"testing"
)

func TestTimestamp(t *testing.T) {
	if len(timestamp()) == 0 {
		t.Errorf("timestamp() returned an empty string")
	}
}

func TestCreatePasswordHash(t *testing.T) {
	h, e := createPasswordHash("test")
	if e != nil {
		t.Errorf("Error returned from createPasswordHash(): " + e.Error())
	}
	if len(h) == 0 {
		t.Errorf("Hash returned from createPasswordHash is empty")
	}
}

func TestCreateUUID(t *testing.T) {
	u, e := createUUID()
	if e != nil {
		t.Errorf("Error returned from createUUID(): " + e.Error())
	}
	if len(u) == 0 {
		t.Errorf("Hash returned from createUUID is empty")
	}
}

func TestCheckHashAgainstPassword(t *testing.T) {
	raw := "test"
	hashed, e := createPasswordHash(raw)
	if e != nil {
		t.Errorf("Error returned from createPasswordHash(): " + e.Error())
	}
	m, err := checkHashAgainstPassword(hashed, raw)
	if err != nil {
		t.Errorf("Error in checkHashAgainstPassword(): " + e.Error())
	}
	if !m {
		t.Errorf("Passwords do not match")
	}
}
