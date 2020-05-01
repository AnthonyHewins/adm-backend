package controllers

import (
	"github.com/gin-gonic/gin"
)

func DcfValuation(c *gin.Context) {
	c.JSON(200, gin.H{"g": "ready for SaaS! add code here."})
}
