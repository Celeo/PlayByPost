package main

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestMain(m *testing.M) {
	m.Run()
	deleteTestData()
}

const envVarDatabasePath = "DB_PATH"
const testDatabaseName = "test_data.db"

func setDBToTest() {
	deleteTestData()
	os.Setenv(envVarDatabasePath, testDatabaseName)
	createTables()
}

func deleteTestData() {
	if _, err := os.Stat(testDatabaseName); os.IsNotExist(err) {
		return
	}
	if err := os.Remove(testDatabaseName); err != nil {
		panic("Could not delete the test database")
	}
}

func addTestUser(db *sqlx.DB) {
	hashedPassword, err := createPasswordHash("password")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(queryCreateUser, "username", hashedPassword, "email")
	if err != nil {
		panic(err)
	}
}

func addTestPost(db *sqlx.DB) {
	_, err := db.Exec(queryCreatePost, "1", timestamp(), "", "New post content")
	if err != nil {
		panic(err)
	}
}
