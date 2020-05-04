package controllers

import (
	"fmt"
	"testing"
	"time"

	"github.com/AnthonyHewins/adm-backend/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestPasswordReset(t *testing.T) {
	db, r := buildRouterAndDB(t)

	test := func (httpcode int, errorcode, msg string, body interface{}) {
		resp, code := buildRequestFn(r, "POST", testPwReset, body)
		assert.Equal(t, httpcode, code)
		assert.Equal(t, errorcode, resp["error"])
		assert.Equal(t, msg, resp["message"])
	}

	// No bind -> 400
	test(400, ERR_PARAM, "invalid request", nil)

	// Can't find user -> 200, don't let them know the email exists
	email := "asiuhdi"
	test(200, "", fmt.Sprintf("if the account for %v exists, an email has been sent to reset the password", email), &pwReset{Email: email})

	// TODO this may be better serviced by re-sending the confirmation link
	// User found, email not confirmed -> do nothing, but tell the user 200 OK
	u := models.User{Email: fmt.Sprintf("asdd%v@sdf.co", time.Now().UnixNano()), Password: "asiuohjdoasjd"}
	u.Create(db)
	test(200, "", fmt.Sprintf("if the account for %v exists, an email has been sent to reset the password", u.Email), &pwReset{Email: u.Email})

	// No reset tokens generated
	upr := models.UserPasswordReset{UserID: u.ID}
	err := db.Where("user_id = ?", u.ID).First(&upr).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	// User exists, email confirmed -> Everything should work
	uec := models.UserEmailConfirmation{UserID: u.ID}
	uec.ConfirmEmail(db)
	test(200, "", fmt.Sprintf("if the account for %v exists, an email has been sent to reset the password", u.Email), &pwReset{Email: u.Email})

	assert.False(t, db.Where("user_id = ?", upr.UserID).First(&upr).RecordNotFound())
}
