package controllers

import (
	"fmt"

	"github.com/AnthonyHewins/adm-backend/models"
	"github.com/gin-gonic/gin"
)

type pwReset struct {
	Email string `json:"email"`
}

func PasswordReset(c *gin.Context) {
	var pwResetForm pwReset
	if !forceBind(c, &pwResetForm) { return }

	db := connectOrError(c)
	if db == nil { return }

	user := models.User{}
	query := db.Where("email = ?", pwResetForm.Email).First(&user)

	if query.RecordNotFound() {
		// Don't let an attacker know they haven't found anything
		fillWithOk(c)
		return
	} else if query.Error != nil {

	}

	upr := models.UserPasswordReset{UserID: user.ID}
	err := upr.ResetPassword(db)
}

func fillWithOk(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("if the account for %v exists, an email has been sent to reset the password.", pwResetForm.Email),
		})
}
