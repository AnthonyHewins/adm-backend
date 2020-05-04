package controllers

import (
	"fmt"
	"testing"
	"time"

	"github.com/AnthonyHewins/adm-backend/models"
	"github.com/stretchr/testify/assert"
)

func TestConfirmPwReset(t *testing.T) {
	db, r := buildRouterAndDB(t)

	test := func(httpcode int, errcode, msg string, body interface{}) {
		resp, code := buildRequestFn(r, "POST", testConfirmPwReset, body)
		assert.Equal(t, httpcode, code)
		assert.Equal(t, errcode, resp["error"])
		assert.Equal(t, msg, resp["message"])
	}

	// Fails binding -> 400 error
	test(400, ERR_PARAM, "invalid request", nil)

	// Bad token -> forbidden
	test(401, ERR_UNAUTHORIZED, "unauthorized", models.UserPasswordReset{Token: "asjdouiasj"})

	// Build the state needed in the DB for the next tests
	u := models.User{Email: fmt.Sprintf("r22r%v@sdmjoif.co", time.Now().UnixNano()), Password: "sndfojus"}
	u.Create(db)

	upr := models.UserPasswordReset{UserID: u.ID}
	upr.CreateResetPasswordToken(db)

	// Token matches, bad pw -> PW too simple
	test(422, ERR_PASSWORD, models.PasswordTooSimple.Error(), &pwResetConfirm{Token: upr.Token, Password: "sdfiu"})

	// Token matches, timed out -> TokenTimeout
	oldTimeoutValue := models.TokenTimeoutThreshold
	models.TokenTimeoutThreshold = 0
	test(403, ERR_LATE, models.TokenTimeout.Error(), &pwResetConfirm{Token: upr.Token, Password: "sdujifhisudfh"})

	models.TokenTimeoutThreshold = oldTimeoutValue

	// Token matches, good pw, in time -> Success
	test(200, "", "password reset successfully", &pwResetConfirm{Token: upr.Token, Password: "asjunbfiusd"})
}
