package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
	"github.com/spf13/viper"
)

func main() {
	loadConfig()
	createTables()
	r := gin.Default()
	r.GET("/", viewIndex)
	r.Run(":5000")
}

func loadConfig() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func createTables() {
	db := database()
	defer db.Close()
	db.MustExec(schemaCreateTables)
}
