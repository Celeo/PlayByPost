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
	h, err := createPasswordHash("test")
	if err != nil {
		t.Errorf("Error returned from createPasswordHash(): " + err.Error())
	}
	if len(h) == 0 {
		t.Errorf("Hash returned from createPasswordHash is empty")
	}
}

func TestCreateUUID(t *testing.T) {
	u, err := createUUID()
	if err != nil {
		t.Errorf("Error returned from createUUID(): " + err.Error())
	}
	if len(u) == 0 {
		t.Errorf("Hash returned from createUUID is empty")
	}
}

func TestCheckHashAgainstPassword(t *testing.T) {
	raw := "test"
	hashed, err := createPasswordHash(raw)
	if err != nil {
		t.Errorf("Error returned from createPasswordHash(): " + err.Error())
	}
	m, err := checkHashAgainstPassword(hashed, raw)
	if err != nil {
		t.Errorf("Error in checkHashAgainstPassword(): " + err.Error())
	}
	if !m {
		t.Errorf("Passwords do not match")
	}
}

func TestInsertRolls(t *testing.T) {
	// TODO
}
