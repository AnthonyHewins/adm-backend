package controllers

import (
	"github.com/AnthonyHewins/adm-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type pwResetConfirm struct {
	Token    string `json:"token"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

func ConfirmPwReset(c *gin.Context) {
	var resetForm pwResetConfirm
	if !forceBind(c, &resetForm) { return }

	db := connectOrError(c)
	if db == nil { return }

	upr := models.UserPasswordReset{Token: resetForm.Token}
	err := upr.ResetPassword(db, resetForm.Password)

	switch err {
	case models.TokenTimeout:
		c.JSON(403, gin.H{
			"error": ERR_LATE,
			"message": err.Error(),
		})
	case gorm.ErrRecordNotFound:
		c.JSON(401, gin.H{
			"error": ERR_UNAUTHORIZED,
			"message": "unauthorized",
		})
	case models.PasswordTooSimple:
		c.JSON(422, gin.H{
			"error": ERR_PASSWORD,
			"message": err.Error(),
		})
	case nil:
		c.JSON(200, gin.H{"message": "password reset successfully"})
	default:
		c.JSON(500, gin.H{"error": ERR_GENERAL, "message": err.Error()})
	}
}
