package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
)

func databaseError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"message": "Could not access database",
		"error":   err.Error(),
	})
}

func timestamp() string {
	return time.Now().UTC().Format("Jan _2, 2006 @ 15:04:05")
}

func createPasswordHash(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), 0)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPasswordAgainstHash(raw, hashed string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
	return err != nil, err
}

func createUUID() (string, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return u4.String(), err
}
