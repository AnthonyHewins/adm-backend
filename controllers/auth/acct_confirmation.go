package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/AnthonyHewins/adm-backend/controllers/api"
	"github.com/AnthonyHewins/adm-backend/models"
	"github.com/AnthonyHewins/adm-backend/smtp"
)

func AcctConfirmation(c *gin.Context) {
	token := c.Query("token")

	// Don't let a botnet know this is a valid link
	if token == "" {
		c.String(404, "page not found")
		return
	}

	api.RequireDB(c, func(db *gorm.DB) {
		uec := findConfirmation(c, db, token)
		if uec == nil { return }

		switch err := uec.ConfirmEmail(db); err {
		case models.TokenTimeout:
			sendNewConfirmationToken(uec.UserID, db, c)
		case nil:
			api.ToAffirmative(c, "email confirmed, welcome")
		default:
			api.Error(c, 500, err.Error())
		}
	})
}

func findConfirmation(c *gin.Context, db *gorm.DB, token string) *models.UserEmailConfirmation {
	uec := models.UserEmailConfirmation{}

	err := db.Where("token = ?", token).First(&uec).Error
	switch err {
	case gorm.ErrRecordNotFound:
		// Don't let an attacker know
		c.String(404, "page not found")
	case nil:
		return &uec
	default:
		api.Error(c, 500, err.Error())
	}

	return nil
}

func sendNewConfirmationToken(id uint64, db *gorm.DB, c *gin.Context) {
	u := models.User{ID: id}
	db.First(&u)

	token, err := u.RefreshConfirmationToken(db)
	if err != nil {
		api.Error(c, 500, err.Error())
		return
	}

	var msg string
	if err = smtp.TokenRefresh(u.Email, token); err == nil {
		msg = "confirmed email too late; another email has been sent"
	} else {
		msg = fmt.Sprintf(
			"confirmed email too late; tried sending another email but failed: %v",
			err.Error(),
		)
	}

	c.JSON(422, api.ToError(ErrLate, msg))
}
