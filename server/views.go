package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func viewIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Index page",
	})
}
