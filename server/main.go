package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// createTables()
	r := gin.Default()
	r.GET("/", viewIndex)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	r.Run(":" + port)
}

func createTables() {
	db := database()
	defer db.Close()
	db.MustExec(schemaCreateTables)
}
