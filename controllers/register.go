package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/AnthonyHewins/adm-backend/models"
)

type registrationForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var form registrationForm

	if !forceBind(c, &form) { return }

	db := connectOrError(c)
	if db == nil { return }
	defer db.Close()

	_, err := models.CreateUser(db, form.Email, form.Password)

	switch err {
	case models.InvalidEmail:
		c.JSON(422, gin.H{
			"error": ERR_EMAIL,
			"message": err.Error(),
		})
	case nil:
		c.JSON(200, gin.H{
			"message": "email registered; please confirm it with the email that was sent to your address.",
		})
	default:
		c.JSON(500, gin.H{
			"error": ERR_GENERAL,
			"message": err.Error(),
		})
	}
}
