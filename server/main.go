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
	r.POST("/login", viewLogin)
	r.POST("/register", viewRegister)
	// TODO apply middleware to the below handlers
	r.GET("/post", viewPosts)
	r.POST("/post", viewCreatePost)
	r.PUT("/post", viewEditPost)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	r.Run(":" + port)
}

func middlewareLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		u := user{}
		s := session{}
		db := database()
		defer db.Close()
		if err := db.Get(&s, querySelectSessionByUUID, header); err != nil {
			databaseError(c, err)
			return
		}
		if err := db.Get(&u, querySelectUserByID, s.UserID); err != nil {
			databaseError(c, err)
			return
		}
		c.Set("authName", u.Name)
		c.Set("authID", u.ID)
		c.Next()
	}
}
