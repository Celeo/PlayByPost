package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // sqolite3 driver
)

func database() *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", "data.db")
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
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	password TEXT NOT NULL,
	email TEXT NOT NULL DEFAULT '',
	notifyOnUpdate INTEGER NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS "post" (
	id INTEGER PRIMARY KEY,
	user_id INTEGER REFERENCES "user" (id) NOT NULL,
	date TEXT NOT NULL,
	content TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "session" (
	user_id INTEGER PRIMARY KEY REFERENCES "user" (id) NOT NULL,
	uuid TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "roll" (
	id INTEGER PRIMARY KEY,
	post_id INTEGER REFERENCES "post" (id) NOT NULL,
	roll TEXT NOT NULL
);
`
const querySelectSessionByUUID = `SELECT * FROM session WHERE uuid=?`
const querySelectUserByID = `SELECT * FROM user WHERE id=?`
const querySelectUserByName = `SELECT * FROM user WHERE name=?`
const querySelectPosts = `SELECT * FROM post`
const querySelectSinglePost = `SELECT * FROM post WHERE id=?`
const queryselectUsers = `SELECT * FROM user`
const queryCreateUser = `INSERT INTO user (name, password, email) VALUES (?, ?, ?)`
const queryCreateSession = `INSERT INTO session VALUES (?, ?)`
const queryDeleteSessionsForUser = `DELETE FROM session WHERE user_id=?`
const queryCreatePost = `INSERT INTO post (user_id, date, content) VALUES (?, ?, ?)`
const queryEditPost = `UPDATE post SET content=? WHERE id=?`
