package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func viewIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Index page",
	})
}

func viewRegister(c *gin.Context) {
	data := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Code     string `json:"code"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		return
	}
	if data.Code != os.Getenv("JOIN_CODE") {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Incorrect join code supplied",
		})
		return
	}
	db := database()
	defer db.Close()
	if _, err := db.Exec(queryCreateUser, data.Name, data.Password, data.Email); err != nil {
		databaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
	})
}

func viewLogin(c *gin.Context) {
	data := struct {
		Name     string `json:"name"`
		Password string `json:"passwword"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		return
	}
	db := database()
	defer db.Close()
	u := user{}
	if err := db.Get(&u, querySelectUserByName, data.Name); err != nil {
		databaseError(c, err)
		return
	}
	uuid, err := createUUID()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error generating UUID for user session",
		})
		return
	}
	if _, err := db.Exec(queryCreateSession, u.ID, uuid); err != nil {
		databaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"uuid":    uuid,
	})
}

func viewPosts(c *gin.Context) {
	posts := []post{}
	db := database()
	defer db.Close()
	if err := db.Select(&posts, querySelectPosts); err != nil {
		databaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, posts)
}

func viewCreatePost(c *gin.Context) {
	data := struct {
		Content string `json:"content"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		return
	}
	db := database()
	defer db.Close()
	if _, err := db.Exec(queryCreatePost, c.GetInt("authID"), timestamp(), data.Content); err != nil {
		databaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Message posted successfully",
	})
}

func viewEditPost(c *gin.Context) {
	data := struct {
		ID      int    `json:"id"`
		Content string `json:"content"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		databaseError(c, err)
		return
	}
	db := database()
	defer db.Close()
	if _, err := db.Exec(queryEditPost, data.ID, data.Content); err != nil {
		databaseError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
