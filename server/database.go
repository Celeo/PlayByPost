package main

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
)

func database() *sqlx.DB {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	return db
}

func createTables() {
	db := database()
	defer db.Close()
	db.MustExec(queryCreateTables)
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

type session struct {
	UserID int    `db:"user_id" json:"userID"`
	UUID   string `db:"uuid" json:"uuid"`
}

const queryCreateTables = `
CREATE TABLE IF NOT EXISTS "user" (
	id bigserial PRIMARY KEY,
	name varchar NOT NULL,
	password varchar NOT NULL,
	email varchar,
	notifyOnUpdate boolean NOT NULL DEFAULT false
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
const querySelectSessionByUUID = `SELECT * FROM "session" WHERE uuid=?`
const querySelectUserByID = `SELECT * FROM "user" WHERE id=?`
const querySelectUserByName = `SELECT * FROM "user" WHERE name=?`
const querySelectPosts = `SELECT * FROM "post"`
const queryCreateUser = `INSERT INTO "user" (name, password, email) VALUES (?, ?, ?)`
const queryCreateSession = `INSERT INTO "session" VALUES (?, ?)`
const queryCreatePost = `INSERT INTO "post" (user_id, date, content) VALUES (?, ?, ?)`
const queryEditPost = `UPDATE "post" SET content=? WHERE id=?`
