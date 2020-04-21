package models

import (
	"os"
	"fmt"
	"time"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestMain(m *testing.M) {
	DBSetupTest(nil)
	os.Exit(m.Run())
}

func TestCreateUser(t *testing.T) {
	db := getDB(t)

	// PW too short
	_, err := CreateUser(db, "email@email.com", "doisjfouijsaofijaosdifj"[:passwordLength - 1])
	if err != PasswordTooSimple {
		t.Errorf("Should've gotten PW too simple, but got %v", err)
	}

	_, err = CreateUser(db, "invalid email", "password")
	if err != InvalidEmail {
		t.Errorf("Should've gotten InvalidEmail, but got %v", err)
	}

	u, err := CreateUser(db, fmt.Sprintf("v%v@email.com", time.Now().UnixNano()), "validPassword1234")
	if err == nil {
		uec := UserEmailConfirmation{}
		query := db.Where("user_id = ?", u.ID).First(&uec)

		if query.Error == nil {
			if uec.UserID != u.ID || uec.Token == "" {
				t.Errorf("somehow the generated token was invalid: %v", uec)
			}
		} else {
			t.Errorf("should've been able to find token for user registration, didn't: %v", query.Error)
		}
	} else {
		t.Errorf("should've successfully created user, but got %v", err)
	}
}

func TestRefreshConfirmationToken(t *testing.T) {
	db := getDB(t)

	u, _ := CreateUser(db, fmt.Sprintf("realUser%v@sdmf.com", time.Now().UnixNano()), "osaidjfio")

	// Get the newly created confirmation for testing
	uec := UserEmailConfirmation{}
	db.Where("user_id = ?", u.ID).First(&uec)
	db.Where("id = ?",      u.ID).First(&u)

	// Timestamp the moment after creation, before refreshing for a test
	oldRegisterTime := u.RegisteredAt.UnixNano()

	token, _ := u.RefreshConfirmationToken(db)

	if token == uec.Token {
		t.Errorf("token didn't change after refresh: %v == %v", token, uec.Token)
	}

	// Re-fetch for tests
	db.Where("id = ?", u.ID).First(&u)

	if u.RegisteredAt.UnixNano() <= oldRegisterTime {
		t.Errorf(
			"Registered at should have been refreshed to a later time period, but may not have: %v <= %v",
			u.RegisteredAt.UnixNano(),
			oldRegisterTime,
		)
	}

	if u.ConfirmedAt != nil {
		t.Errorf("user was confirmed even though all that happened was a token refresh at date %v", u.ConfirmedAt)
	}

	// Already confirmed email should result in error (after refreshing)
	db.Model(&u).Update("confirmed_at", time.Now())
	db.Where("id = ?", u.ID).First(&u)

	if _, err := u.RefreshConfirmationToken(db); err != AlreadyConfirmed {
		t.Errorf("already confirmed email should've thrown AlreadyConfirmed, got %v", err)
	}
}

func getDB(t *testing.T) *gorm.DB {
	db, err := Connect()
	if err != nil { t.Fatal(err.Error()) }

	return db
}
