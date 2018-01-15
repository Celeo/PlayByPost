package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func database() *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", viper.GetString("database.path"))
	if err != nil {
		panic("Could not open database")
	}
	return db
}

const schemaCreateTables = `
CREATE TABLE IF NOT EXISTS user (

);

CREATE TABLE IF NOT EXISTS post (

);
`
