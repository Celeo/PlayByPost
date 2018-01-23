package main

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
)

func database() *sqlx.DB {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "user=postgres password=postgres host=localhost port=5432 dbname=playbypost sslmode=disable"
	}
	db, err := sqlx.Connect("postgres", databaseURL)
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

// A User is a member
type User struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Password       string `json:"-"`
	Email          string `json:"email"`
	NotifyOnUpdate bool   `json:"notifyOnUpdate"`
}

// A Post is a message
type Post struct {
	ID      int    `json:"id"`
	UserID  int    `db:"user_id" json:"userID"`
	Date    string `json:"date"`
	Content string `json:"content"`
}

// A Session is used for authentication
type Session struct {
	UserID int    `db:"user_id" json:"userID"`
	UUID   string `db:"uuid" json:"uuid"`
}

// A Roll is the result of rolling one or more dice
type Roll struct {
	ID     int    `json:"id"`
	PostID int    `db:"post_id" json:"postID"`
	Roll   string `json:"roll"`
}

const queryCreateTables = `
CREATE TABLE IF NOT EXISTS "user" (
	id bigserial PRIMARY KEY,
	name varchar NOT NULL,
	password varchar NOT NULL,
	email varchar NOT NULL DEFAULT '',
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

CREATE TABLE IF NOT EXISTS "roll" (
	id bigserial PRIMARY KEY,
	post_id bigserial REFERENCES "post" (id) NOT NULL,
	roll varchar NOT NULL
);
`
const querySelectSessionByUUID = `SELECT * FROM "session" WHERE uuid=$1`
const querySelectUserByID = `SELECT * FROM "user" WHERE id=$1`
const querySelectUserByName = `SELECT * FROM "user" WHERE name=$1`
const querySelectPosts = `SELECT * FROM "post"`
const querySelectSinglePost = `SELECT * FROM "post" WHERE id=$1`
const queryselectUsers = `SELECT * FROM "user"`
const queryCreateUser = `INSERT INTO "user" (name, password, email) VALUES ($1, $2, $3)`
const queryCreateSession = `INSERT INTO "session" VALUES ($1, $2)`
const queryDeleteSessionsForUser = `DELETE FROM "session" WHERE user_id=$1`
const queryCreatePost = `INSERT INTO "post" (user_id, date, content) VALUES ($1, $2, $3)`
const queryEditPost = `UPDATE "post" SET content=? WHERE id=$1`
