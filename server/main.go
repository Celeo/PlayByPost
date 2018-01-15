package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	createTables()
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

func middlewareLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}
		// TODO
		c.Next()
	}
}
