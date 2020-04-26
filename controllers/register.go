package controllers

import (
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/AnthonyHewins/adm-backend/models"
)

var (
	uniquenessViolation = regexp.MustCompile("duplicate key value violates unique constraint")
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

	_, err := models.CreateUser(db, form.Email, form.Password)
	db.Close()

	switch err {
	case models.InvalidEmail:
		c.JSON(422, gin.H{
			"error": ERR_EMAIL,
			"message": err.Error(),
		})
	case models.PasswordTooSimple:
		c.JSON(422, gin.H{
			"error": ERR_PASSWORD,
			"message": err.Error(),
		})
	case nil:
		c.JSON(200, gin.H{
			"message": "email registered; please confirm it with the email that was sent to your address.",
		})
	default:
		errstring := err.Error()

		if uniquenessViolation.MatchString(errstring) {
			c.JSON(422, gin.H{
				"error": ERR_ALREADY_EXISTS,
				"message": "email is taken, please use a different one",
			})
		} else {
			c.JSON(500, gin.H{
				"error": ERR_GENERAL,
				"message": errstring,
			})
		}
	}
}
