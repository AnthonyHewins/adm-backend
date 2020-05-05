package auth

import (
	"fmt"

	"github.com/AnthonyHewins/adm-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
	err := db.Where("email = ?", pwResetForm.Email).First(&user).Error

	switch err {
	case gorm.ErrRecordNotFound:
		// Don't let an attacker know they found an email
		fillWithOk(c, pwResetForm.Email)
	case nil:
		performReset(c, db, &user)
	default:
		c.JSON(500, gin.H{
			"error": ERR_GENERAL,
			"message": err.Error(),
		})
	}

}

func performReset(c *gin.Context, db *gorm.DB, user *models.User) {
	// Email not confirmed -> they can't reset it
	if user.ConfirmedAt == nil {
		fillWithOk(c, user.Email)
		return
	}

	upr := models.UserPasswordReset{UserID: user.ID}
	if err := upr.CreateResetPasswordToken(db); err != nil {
		c.JSON(500, gin.H{
			"error": ERR_GENERAL,
			"message": err.Error(),
		})
	}

	fillWithOk(c, user.Email)
}

func fillWithOk(c *gin.Context, email string) {
	c.JSON(200, gin.H{
		"message": fmt.Sprintf("if the account for %v exists, an email has been sent to reset the password", email),
	})
}
