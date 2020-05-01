package controllers

import (
	"fmt"
	"testing"
	"time"

	"github.com/AnthonyHewins/adm-backend/models"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/stretchr/testify/assert"
)

func TestUnauthorizedHandler(t *testing.T) {
	_, router := buildRouterAndDB(t)

	// Test no token is all that's needed; the library should handle the rest.
	// We just need to see that the response is 401 (since no auth was provided),
	// and that it has ERR_UNAUTHORIZED.
	resp, code := buildRequestFn(router, "GET", testTicker, nil)
	assert.Equal(t, 401, code)
	assert.Equal(t, ERR_UNAUTHORIZED, resp["error"])
}

func TestAuthenticate(t *testing.T) {
	db, router := buildRouterAndDB(t)

	test := func(c int, error, message string, body interface{}) {
		resp, code := buildRequestFn(router, "POST", testLogin, body)
		assert.Equal(t, c, code)
		assert.Equal(t, error, resp["error"])
		if message != "" {
			assert.Equal(t, message, resp["message"])
		}
	}

	// Missing login values -> unauthorized
	test(401, ERR_UNAUTHORIZED, jwt.ErrMissingLoginValues.Error(), nil)

	// Nonexisting user -> unauthorized
	test(401, ERR_UNAUTHORIZED, jwt.ErrFailedAuthentication.Error(), &login{Email: "1", Password: "uiodsjdsg"})

	// User exists, pw correct, but not confirmed -> unauthorized
	pw := "ihjfoiurdes"
	u, _ := models.CreateUser(db, fmt.Sprintf("ij%v@dfk.com", time.Now().UnixNano()), pw)
	login := login{ Email: u.Email, Password: pw }
	test(401, ERR_UNAUTHORIZED, "you have not confirmed your email yet; please confirm it to log in", &login)

	// User email confirmed, but password doesn't match -> unauthorized
	uec := models.UserEmailConfirmation{}
	db.Where("user_id = ?", u.ID).First(&uec)
	uec.ConfirmEmail(db)

	login.Password = "111"
	test(401, ERR_UNAUTHORIZED, "", &login)

	// User email confirmed, PW matches -> successful login
	login.Password = pw
	test(200, "", "", &login)
}

/*
TODO may be used in the future, more likely than not

func TestAuthorizator(t *testing.T) {

}
/*func TestLogoutHandler(t *testing.T) {

}
*/
