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
	joinCode := os.Getenv("JOIN_CODE")
	if joinCode == "" {
		joinCode = "join"
	}
	if data.Code != joinCode {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Incorrect join code supplied",
		})
		return
	}
	db := database()
	defer db.Close()
	hashedPassword, err := createPasswordHash(data.Password)
	if err != nil {
		abortError(c, err)
		return
	}
	_, err = db.Exec(queryCreateUser, data.Name, hashedPassword, data.Email)
	if err != nil {
		abortError(c, err)
		return
	}
	u, err := getUserByName(data.Name)
	if err != nil {
		abortError(c, err)
		return
	}
	uuid, err := createSession(u)
	if err != nil {
		abortError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
		"uuid":    uuid,
		"name":    data.Name,
	})
}

func viewLogin(c *gin.Context) {
	data := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		return
	}
	db := database()
	defer db.Close()
	u, err := getUserByName(data.Name)
	if err != nil {
		abortError(c, err)
		return
	}
	passwordMatch, err := checkHashAgainstPassword(u.Password, data.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error checking hashed password for login",
			"error":   err.Error(),
		})
		return
	}
	if !passwordMatch {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Password mismatch",
		})
		return
	}
	uuid, err := createSession(u)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error generating UUID for user session",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"uuid":    uuid,
		"name":    u.Name,
	})
}

func viewPosts(c *gin.Context) {
	posts := []post{}
	users := []user{}
	userMap := make(map[int]user)
	type returnPost struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Date    string `json:"date"`
		Content string `json:"content"`
	}
	retVal := []returnPost{}
	db := database()
	defer db.Close()
	if err := db.Select(&posts, querySelectPosts); err != nil {
		abortError(c, err)
		return
	}
	if err := db.Select(&users, queryselectUsers); err != nil {
		abortError(c, err)
		return
	}
	for _, u := range users {
		userMap[u.ID] = u
	}
	for _, p := range posts {
		retVal = append(retVal, returnPost{
			ID:      p.ID,
			Name:    userMap[p.UserID].Name,
			Date:    p.Date,
			Content: p.Content,
		})
	}
	c.JSON(http.StatusOK, retVal)
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
		abortError(c, err)
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
		abortError(c, err)
		return
	}
	db := database()
	defer db.Close()
	if _, err := db.Exec(queryEditPost, data.ID, data.Content); err != nil {
		abortError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
