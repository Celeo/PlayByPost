package main

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // sqolite3 driver
)

func database() *sqlx.DB {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data.db"
	}
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	return db
}

// createTables executes the query to creates all of the database tables.
func createTables() {
	db := database()
	defer db.Close()
	db.MustExec(queryCreateTables)
}

// A User is a member
type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Password     string `json:"-"`
	Email        string `json:"email"`
	PostsPerPage int    `json:"postsPerPage" db:"postsPerPage"`
	NewestAtTop  bool   `json:"newestAtTop" db:"newestAtTop"`
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
	ID     int    `json:"id"`
	UserID int    `db:"user_id" json:"userID"`
	UUID   string `db:"uuid" json:"uuid"`
}

// A Roll is the result of rolling some dice
type Roll struct {
	ID      int    `json:"id"`
	UserID  int    `db:"user_id" json:"userID"`
	PostID  int    `db:"post_id" json:"postID"`
	Pending bool   `json:"pending"`
	String  string `json:"string"`
	Value   int    `json:"value"`
}

const queryCreateTables = `
CREATE TABLE IF NOT EXISTS "user" (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	password TEXT NOT NULL,
	email TEXT NOT NULL DEFAULT '',
	postsPerPage INT NOT NULL DEFAULT 25,
	newestAtTop INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "post" (
	id INTEGER PRIMARY KEY,
	user_id INTEGER REFERENCES "user" (id) NOT NULL,
	date TEXT NOT NULL,
	content TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "session" (
	id INTEGER PRIMARY KEY,
	user_id INTEGER REFERENCES "user" (id) NOT NULL,
	uuid TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "roll" (
	id INTEGER PRIMARY KEY,
	user_id INTEGER REFERENCES "user" (id) NOT NULL,
	post_id INTEGER REFERENCES "post" (id) NOT NULL DEFAULT 0,
	pending INT NOT NULL DEFAULT 1,
	string TEXT NOT NULL DEFAULT '',
	value INTEGER NOT NULL DEFAULT 0
);
`
const querySelectSessionByUUID string = `SELECT * FROM session WHERE uuid=?`
const querySelectUserByID string = `SELECT * FROM user WHERE id=?`
const querySelectUserByName string = `SELECT * FROM user WHERE name=?`
const querySelectPosts string = `SELECT * FROM post`
const querySelectSinglePost string = `SELECT * FROM post WHERE id=?`
const querySelectUsers string = `SELECT * FROM user`
const queryCreateUser string = `INSERT INTO user (name, password, email) VALUES (?, ?, ?)`
const queryCreateSession string = `INSERT INTO session (user_id, uuid) VALUES (?, ?)`
const queryDeleteSessionsForUser string = `DELETE FROM session WHERE user_id=?`
const queryCreatePost string = `INSERT INTO post (user_id, date, content) VALUES (?, ?, ?)`
const queryEditPost string = `UPDATE post SET content=? WHERE id=?`
const queryUpdatePassword string = `UPDATE user SET password=? WHERE id=?`
const queryUpdateUser string = `UPDATE user SET name=?, email=?, postsPerPage=?, newestAtTop=? WHERE id=?`
const queryGetPendingRollsForUser string = `SELECT * FROM roll WHERE user_id=? AND pending=1`
const queryInsertPendingRoll string = `INSERT INTO roll (user_id, string, value) VALUES (?, ?, ?)`
const querySavePendingRoll string = `UPDATE roll SET pending=0, post_id=? WHERE user_id=? AND pending=1`
const querySelectSavedRolls string = `SELECT * FROM roll WHERE pending=0`
const queryInvalidLogins string = `DELETE FROM session WHERE user_id=? AND uuid!=?`
const queryGetAllSessions string = `SELECT * FROM session`
