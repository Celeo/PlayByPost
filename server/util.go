package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
)

func abortError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"message": "Could not access database",
		"error":   err.Error(),
	})
}

func timestamp() string {
	return time.Now().UTC().Format("Jan _2, 2006 @ 15:04:05")
}

func createPasswordHash(raw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), 0)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func checkHashAgainstPassword(hashed, raw string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
	return err == nil, err
}

func createUUID() (string, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return u4.String(), err
}

func createSession(u user) (string, error) {
	db := database()
	defer db.Close()
	uuid, err := createUUID()
	if err != nil {
		return "", err
	}
	if _, err := db.Exec(queryDeleteSessionsForUser, u.ID); err != nil {
		return "", err
	}
	if _, err := db.Exec(queryCreateSession, u.ID, uuid); err != nil {
		return "", err
	}
	return uuid, nil
}
