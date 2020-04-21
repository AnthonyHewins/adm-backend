package controllers

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
	router := buildRouter()

	models.DBSetupTest(nil) // nil => use default configuration
	db, err := models.Connect()
	if err != nil { t.Fatalf("test failed because test DB broke: %v", err) }

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

	test := func(path string, eCode int, err string)  {
		resp, code := buildRequestFn(router, "GET", path, nil)
		assert.Equal(t, eCode, code)
		assert.Equal(t, err, resp["error"])
	}

	// registered 15 min late, the cutoff time
	u := models.User{
		Email: fmt.Sprintf("fake%v@gmail.com", time.Now()),
		Password: "iasdjoaisd",
		RegisteredAt: time.Now().Add(time.Hour),
	}

	if err := db.Create(&u).Error; err != nil { t.Fatal(err.Error()) }

	uec := models.UserEmailConfirmation{UserID: u.ID}
	if err := db.Model(uec).Related(&u).Error; err != nil { t.Fatal(err.Error()) }

	test(url(uec.Token), 422, ERR_LATE)

	// On time
	u.RegisteredAt = time.Now()
	db.Save(&u)
	test(url(uec.Token), 200, "")
}

func url(token string) string {
	return fmt.Sprintf("%v?token=%v", testConfirmation, token)
}
