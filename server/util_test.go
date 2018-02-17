package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTimestamp(t *testing.T) {
	require.NotEmpty(t, timestamp())
}

func TestCreatePasswordHash(t *testing.T) {
	h, err := createPasswordHash("test")
	require.Nil(t, err)
	require.NotEmpty(t, h)
}

func TestCreateUUID(t *testing.T) {
	h, err := createUUID()
	require.Nil(t, err)
	require.NotEmpty(t, h)
}

func TestCheckHashAgainstPassword(t *testing.T) {
	raw := "test"
	hashed, err := createPasswordHash(raw)
	require.Nil(t, err)
	m := checkHashAgainstPassword(hashed, raw)
	require.True(t, m)
}

func TestGetUserByID(t *testing.T) {
	newDB()
	db := database()
	defer db.Close()
	addUser(db)
	u, err := getUserByID(1)
	require.Nil(t, err)
	require.NotNil(t, u)
	require.Equal(t, u.ID, 1)
}

func TestGetUserByIDNoUser(t *testing.T) {
	newDB()
	u, err := getUserByID(1)
	require.NotNil(t, err)
	require.NotNil(t, u)
	require.Equal(t, u.ID, 0)
}

func TestGetUserByName(t *testing.T) {
	newDB()
	db := database()
	defer db.Close()
	addUser(db)
	u, err := getUserByName("username")
	require.Nil(t, err)
	require.NotNil(t, u)
	require.Equal(t, u.ID, 1)
}

func TestGetUserByNameNoUser(t *testing.T) {
	newDB()
	u, err := getUserByName("username")
	require.NotNil(t, err)
	require.NotNil(t, u)
	require.Equal(t, u.ID, 0)
}
