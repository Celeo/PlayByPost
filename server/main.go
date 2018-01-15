package main

import (
	"net/http"
	"os"

	"github.com/itsjamie/gin-cors"

	"github.com/gin-gonic/gin"
)

func main() {
	createTables()
	r := gin.Default()
	r.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "GET, PUT, POST, DELETE",
		RequestHeaders: "Origin, Authorization, Content-Type",
		Credentials:    true,
	}))

	r.GET("/", viewIndex)
	r.POST("/login", viewLogin)
	r.POST("/register", viewRegister)
	r.GET("/post", middlewareLoggedIn(), viewPosts)
	r.POST("/post", middlewareLoggedIn(), viewCreatePost)
	r.PUT("/post", middlewareLoggedIn(), viewEditPost)

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
			abortError(c, err)
			return
		}
		if err := db.Get(&u, querySelectUserByID, s.UserID); err != nil {
			abortError(c, err)
			return
		}
		c.Set("authName", u.Name)
		c.Set("authID", u.ID)
		c.Next()
	}
}
