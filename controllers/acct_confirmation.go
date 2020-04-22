package controllers

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"

	"github.com/AnthonyHewins/adm-backend/models"
	"github.com/AnthonyHewins/adm-backend/smtp"
)

func AcctConfirmation(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		c.JSON(404, gin.H{"message":"page not found"})
		return
	}

	db := connectOrError(c)
	if db == nil { return }
	defer db.Close()

	uec := findConfirmation(c, db, token)
	if uec == nil { return }

	switch err := uec.ConfirmEmail(db); err {
	case &models.EmailConfirmationLate:
		sendNewConfirmationToken(uec.UserID, db, c)
	case nil:
		c.JSON(200, gin.H{"message": "email confirmed, welcome"})
	default:
		c.JSON(500, gin.H{
			"error": ERR_GENERAL,
			"message": err.Error(),
		})
	}
}

func findConfirmation(c *gin.Context, db *gorm.DB, token string) *models.UserEmailConfirmation {
	uec := models.UserEmailConfirmation{}
	tokenQuery := db.Where("token = ?", token).First(&uec)

	if tokenQuery.RecordNotFound() {
		c.JSON(404, gin.H{"message":"page not found"})
		return nil
	}

	if tokenQuery.Error != nil {
		c.JSON(500, gin.H{
			"error": ERR_GENERAL,
			"message": tokenQuery.Error.Error(),
		})
		return nil
	}

	return &uec
}

func sendNewConfirmationToken(id uint64, db *gorm.DB, c *gin.Context) {
	u := models.User{ID: id}
	db.First(&u)
	token, err := u.RefreshConfirmationToken(db)
	if err != nil {
		c.JSON(500, gin.H{
			"error": ERR_GENERAL,
			"message": err.Error(),
		})
		return
	}

	if err = smtp.AccountConfirmation(u.Email, token); err != nil {
		c.JSON(422, gin.H{
			"error": ERR_LATE,
			"message": "confirmed email too late; another email has been sent",
		})
	} else {
		c.JSON(500, gin.H{
			"error": ERR_LATE,
			"message": fmt.Sprintf(
				"confirmed email too late; tried sending another email but failed: %v",
				err.Error(),
			),
		})
	}
}
