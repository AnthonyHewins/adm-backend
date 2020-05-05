package auth

import (
	"fmt"
	"time"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/AnthonyHewins/adm-backend/models"
)

func TestAcctConfirmation(t *testing.T) {
	db, router := buildRouterAndDB(t)

	// 404s should occur as a security measure
	test404s(t, router)

	// Everything else
	testConfirm(t, router, db)
}

func test404s(t *testing.T, router *gin.Engine) {
	test := func(path string, eCode int) {
		_, code := buildRequestFn(router, "GET", path, nil)
		assert.Equal(t, eCode, code)
	}

	test(testConfirmation, 404)
	test(url("DNE"),       404)
}

func testConfirm(t *testing.T, router *gin.Engine, db *gorm.DB) {
	test := func(path string, eCode int, err string) {
		resp, code := buildRequestFn(router, "GET", path, nil)
		assert.Equal(t, eCode, code)
		assert.Equal(t, err,   resp["error"])
	}

	// Setup: user is saved with
	uec  := models.UserEmailConfirmation{}
	u := models.User{Email: fmt.Sprintf("fake%v@gmail.com", time.Now().UnixNano()), Password: "iasdjoaisd"}
	u.Create(db)

	oldTimeoutValue := models.TokenTimeoutThreshold

	// Try confirm with failure
	models.TokenTimeoutThreshold = 0
	db.Where("user_id = ?", u.ID).First(&uec)
	test(url(uec.Token), 422, ERR_LATE)

	oldToken := uec.Token
	db.Where("user_id = ?", u.ID).First(&uec)

	if oldToken == uec.Token {
		t.Errorf("token should be new but isn't: (old, new) => %v, %v", oldToken, uec.Token)
	}


	// Try confirm with success
	models.TokenTimeoutThreshold = oldTimeoutValue
	db.Where("user_id = ?", u.ID).First(&uec)
	test(url(uec.Token), 200, "")
	db.First(&u)

	if u.ConfirmedAt == nil {
		t.Errorf("user should have been confirmed but wasn't: %v", u)
	}

	shouldntBeFound := db.Where("user_id = ?", u.ID).First(&uec)
	if !shouldntBeFound.RecordNotFound() {
		t.Errorf("all tokens should have been deleted after confirmation but wasn't: %v", shouldntBeFound)
	}
}

func url(token string) string {
	return fmt.Sprintf("%v?token=%v", testConfirmation, token)
}
