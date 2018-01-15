package main

import (
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // sqolite3 driver
	"github.com/spf13/viper"
)

func main() {
	loadConfig()
	// createTables()
	r := gin.Default()
	r.GET("/", viewIndex)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	r.Run(":" + port)
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
