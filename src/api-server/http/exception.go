package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GlobalExceptionWare(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": "500", "message": err.(string)})
			c.Abort()
		}
	}()
	c.Next()
}
