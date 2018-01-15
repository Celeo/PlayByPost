package main

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
)

func database() *sqlx.DB {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	return db
}

type user struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"-"`
}

type post struct {
	ID      int    `json:"id"`
	UserID  int    `db:"user_id" json:"userID"`
	Date    string `json:"date"`
	Content string `json:"content"`
}

const schemaCreateTables = `
CREATE TABLE IF NOT EXISTS "user" (
	id bigserial PRIMARY KEY,
	name varchar NOT NULL,
	password varchar NOT NULL
);

CREATE TABLE IF NOT EXISTS "post" (
	id bigserial PRIMARY KEY,
	user_id bigserial REFERENCES "user" (id) NOT NULL,
	date varchar NOT NULL,
	content varchar NOT NULL
);

CREATE TABLE IF NOT EXISTS "session" (
	user_id bigserial PRIMARY KEY REFERENCES "user" (id) NOT NULL,
	uuid char(36) NOT NULL
);
`

func createPasswordHash(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), 0)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPasswordAgainstHash(raw, hashed string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
	return err != nil, err
}

func createUUID() (string, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return u4.String(), err
}
