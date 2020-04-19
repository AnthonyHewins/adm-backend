package controllers

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"

	"github.com/AnthonyHewins/adm-backend/models"
)

func AcctConfirmation(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		c.JSON(404, gin.H{"message":"page not found"})
		return
	}

	uec := models.UserEmailConfirmation{Token: token}

	db, err := models.Connect()
	defer db.Close()

	if err != nil {
		c.JSON(500, gin.H{
			"err": ERR_GENERAL,
			"message": err.Error(),
		})
	}

	tokenQuery := db.Where("token = ?", token).First(&uec)
	if tokenQuery.RecordNotFound() {
		c.JSON(404, gin.H{"message":"page not found"})
		return
	} else if tokenQuery.Error != nil {
		c.JSON(500, gin.H{
			"err": ERR_GENERAL,
			"message": err.Error(),
		})
		return
	}

	switch err := uec.ConfirmEmail(db); err {
	case &models.EmailConfirmationLate:
		sendNewConfirmationToken(uec.UserId, db, c)
	case nil:
		c.JSON(200, gin.H{"message": "email confirmed, welcome"})
	default:
		c.JSON(500, gin.H{
			"err": ERR_GENERAL,
			"message": err.Error(),
		})
	}
}

func sendNewConfirmationToken(id int64, db *gorm.DB, c *gin.Context) {
	u := models.User{ID: id}
	db.First(&u)
	err := u.RefreshConfirmationToken(db)

	if err == nil {
		c.JSON(422, gin.H{
			"err": ERR_LATE,
			"message": "confirmed email too late; another email has been sent",
		})
	} else {
		c.JSON(500, gin.H{
			"err": ERR_LATE,
			"message": fmt.Sprintf(
				"confirmed email too late; tried sending another email but failed: %v",
				err.Error(),
			),
		})
	}
}
