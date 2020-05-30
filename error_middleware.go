package main

import (
	"github.com/gin-gonic/gin"
)

// ErrorHandler add to the response body the error to be handled in the client-side.
func ErrorHandler(c *gin.Context) {
	c.Next()
	err := c.Errors.Last()
	if err == nil {
		return
	}
	c.IndentedJSON(c.Writer.Status(), gin.H{
		"error": err.Error()})
	c.Abort()

	return
}
