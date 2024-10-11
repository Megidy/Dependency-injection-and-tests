package utils

import "github.com/gin-gonic/gin"

func HandleError(c *gin.Context, err error, msg string, statusCode int) {
	c.JSON(statusCode, gin.H{
		"error":   err,
		"details": msg,
	})

}
