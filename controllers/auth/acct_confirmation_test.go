package auth

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/AnthonyHewins/adm-backend/controllers/api"
	"github.com/AnthonyHewins/adm-backend/models"
)

func TestAcctConfirmation(t *testing.T) {
	db := dbInstance()

	test(
		t, db,
		"",
		nil,
		&api.Error{Http: 404, Code: api.ErrNotFound, Msg: "the server can't find what you're looking for"},
	)

	// Setup: user is saved with
	u := models.User{Email: fmt.Sprintf("fak1e%v@gmail.com", time.Now().UnixNano()), Password: "iasdjoaisd"}
	u.Create(db)

	oldTimeoutValue := models.TokenTimeoutThreshold

	// Try confirm with failure
	models.TokenTimeoutThreshold = 0
	uec  := models.UserEmailConfirmation{}
	db.Where("user_id = ?", u.ID).First(&uec)

	test(
		t, db,
		uec.Token,
		nil,
		&api.Error{Http: 403, Code: ErrLate, Msg: "confirmed email too late; another email has been sent"},
	)

	oldToken := uec.Token
	db.Where("user_id = ?", u.ID).First(&uec)

	if oldToken == uec.Token {
		t.Errorf("token should be new but isn't: (old, new) => %v, %v", oldToken, uec.Token)
	}

	// Try confirm with success
	models.TokenTimeoutThreshold = oldTimeoutValue
	db.Where("user_id = ?", u.ID).First(&uec)

	test(
		t, db,
		uec.Token,
		&api.Affirmative{Msg: "email confirmed, welcome"},
		nil,
	)

	db.First(&u)

	if u.ConfirmedAt == nil {
		t.Errorf("user should have been confirmed but wasn't: %v", u)
	}

	shouldntBeFound := db.Where("user_id = ?", u.ID).First(&uec)
	if !shouldntBeFound.RecordNotFound() {
		t.Errorf("all tokens should have been deleted after confirmation but wasn't: %v", shouldntBeFound)
	}
}

func test(t *testing.T, db *gorm.DB, token string, want1 api.Payload, want2 *api.Error) {
	t.Run(token, func(t *testing.T) {
		got, got1 := acctConfirmation(db, token)
		if !reflect.DeepEqual(got, want1) {
			t.Errorf("got = %v, want %v", got, want1)
		}
		if !reflect.DeepEqual(got1, want2) {
			t.Errorf("got1 = %v, want %v", got1, want2)
		}
	})
}
