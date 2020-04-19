package controllers

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/AnthonyHewins/adm-backend/models"
)

const (
	emailRegexString = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

var emailRegex = regexp.MustCompile(emailRegexString)

type registrationForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var form registrationForm

	if err := c.BindJSON(&form); err != nil {
		c.JSON(400, gin.H{
			"error": ERR_PARAM,
			"message": err.Error(),
		})
		return
	}

	if !emailRegex.MatchString(form.Email) {
		c.JSON(422, gin.H{
			"error": ERR_EMAIL,
			"message": "email doesn't pass validation; is it a valid email?",
		})
		return
	}

	db, err := models.Connect()
	if err != nil {
		c.JSON(500, gin.H{
			"error": ERR_GENERAL,
			"message": err.Error(),
		})
		return
	}

	_, err = models.CreateUser(db, form.Email, form.Password)
	db.Close()

	if err != nil {
		c.JSON(500, gin.H{
			"error": ERR_GENERAL,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "email registered; please confirm it with the email that was sent to your address.",
	})
}
