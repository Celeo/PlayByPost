package main

import (
	"os"

	"github.com/jmoiron/sqlx"
)

func database() *sqlx.DB {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE"))
	if err != nil {
		panic("Could not open database")
	}
	return db
}

const schemaCreateTables = `
CREATE TABLE IF NOT EXISTS user (
	id bigserial PRIMARY KEY,
	name varchar,
	password varchar
);

CREATE TABLE IF NOT EXISTS post (
	id bigserial PRIMARY KEY,
	user_id bigserial,
	date varchar,
	content varchar
);
`
