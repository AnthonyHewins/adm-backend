package api

import (
	"github.com/gin-gonic/gin"
)

func ToAffirmative(c *gin.Context, s string) {
	c.JSON(200, gin.H{
		"message": s,
	})
}
