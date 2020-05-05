package api

import (
	"github.com/AnthonyHewins/adm-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func RequireBind(c *gin.Context, structToBind interface{}, fn func()) {
	if err := c.BindJSON(&structToBind); err != nil {
		Error(c, 400, err.Error())
		return
	}

	fn()
}

func RequireDB(c *gin.Context, fn func(db *gorm.DB)) {
	db, err := models.Connect()

	if err != nil {
		Error(c, 500, err.Error())
		return
	}
	defer db.Close()

	fn(db)
}

func RequireBindAndDB(c *gin.Context, structToBind interface{}, fn func(db *gorm.DB)) {
	if err := c.BindJSON(&structToBind); err != nil {
		Error(c, 400, err.Error())
		return
	}

	db, err := models.Connect()

	if err != nil {
		Error(c, 500, err.Error())
		return
	}

	fn(db)
}
