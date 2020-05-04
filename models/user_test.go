package models

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	DBSetupTest(nil)
	os.Exit(m.Run())
}

func TestCreateUser(t *testing.T) {
	db := getDB(t)

	assertUser := func(err error, user *User) {
		actual := user.Create(db)
		assert.Equal(t, err, actual)
	}

	// Fails
	assertUser(PasswordTooSimple, &User{Email: "email@email.com", Password: "sdfsf"})
	assertUser(InvalidEmail, &User{Email: "invalid email", Password: "password"})

	// Success
	email := fmt.Sprintf("HASCAPS%v@email.com", time.Now().UnixNano())
	u := &User{Email: email, Password: "validPassword1234"}
	assertUser(nil, u)

	// After save -> email is lowercase
	assert.Equal(t, strings.ToLower(email), u.Email)

	// After save -> Password field is gone for security
	assert.Equal(t, "", u.Password)

	// After save -> UserEmailConfirmation exists
	uec := UserEmailConfirmation{}
	query := db.Where("user_id = ?", u.ID).First(&uec)

	if query.Error == nil {
		if uec.UserID != u.ID || uec.Token == "" {
			t.Errorf("somehow the generated token was invalid: %v", uec)
		}
	} else {
		t.Errorf("should've been able to find token for user registration, didn't: %v", query.Error)
	}
}

func TestRefreshConfirmationToken(t *testing.T) {
	db := getDB(t)

	u := createUser(db)

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

func TestAuthenticate(t *testing.T) {
	db := getDB(t)

	// User not found by ID -> RecordNotFound
	u := User{ID: 99999999999 }
	assert.True(t, gorm.IsRecordNotFoundError(u.Authenticate(db)))

	// User not found by email -> RecordNotFound
	u = User{Email: "ansuido"}
	assert.True(t, gorm.IsRecordNotFoundError(u.Authenticate(db)))

	// Create a user
	pw := "sdunfjsd"
	email := fmt.Sprintf("r7rr%v@jidfos.oc", time.Now().UnixNano())
	u = User{Email: email, Password: pw}
	u.Create(db)

	// Wrong PW, email not confirmed -> bcrypt error
	u.Password = "saidnfhiosdjfiosdjpffff"
	assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, u.Authenticate(db))

	// Right PW, email not confirmed (results in same error) -> EmailNotConfirmed
	u.Password = pw
	assert.Equal(t, EmailNotConfirmed, u.Authenticate(db))

	// Confirm the password, refresh the object
	uec := UserEmailConfirmation{UserID: u.ID}
	uec.ConfirmEmail(db)
	db.First(&u) // fetch confirmed at info, currently doesn't have it

	// Wrong PW, email confirmed -> bcrypt mismatch error
	u.Password = "asnduikjsaf"
	assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, u.Authenticate(db))

	// Right PW, email confirmed -> error is nil
	u.Password = pw
	assert.Equal(t, nil, u.Authenticate(db))
}

func TestResetPassword(t *testing.T) {
	db := getDB(t)

	tooShort := "sadd"
	good := "ahfiusdnhfisdfnoi33"
	u := createConfirmedUser(db)

	// PW too short
	u.Password = tooShort
	assert.Equal(t, PasswordTooSimple, u.ResetPassword(db))

	// PW passed validation -> password is now different
	u.Password = good
	u.ResetPassword(db)

	u.Password = good
	assert.Equal(t, nil, u.Authenticate(db))

	// PW passed validation -> password field is nil'd out for security
	u.Password = good
	u.ResetPassword(db)
	assert.Equal(t, "", u.Password)

	// PW passed validation -> all UserPasswordResets are gone
	upr := UserPasswordReset{}
	err := db.Where("user_id = ?", u.ID).First(&upr).Error
	if !gorm.IsRecordNotFoundError(err) {
		t.Errorf("all UserPasswordResets should be gone after the PW has been reset, got this error instead of RecordNotFound: %v", err)
	}
}

func getDB(t *testing.T) *gorm.DB {
	db, err := Connect()
	if err != nil { t.Fatal(err.Error()) }

	return db
}

func createUser(db *gorm.DB) *User {
	u := &User{
		Email: fmt.Sprintf("sdsuhb%v@asdji.co", time.Now().UnixNano()),
		Password: "adsjfasdfasdfa",
	}
	err := u.Create(db)
	if err != nil {
		panic(err)
	}
	return u
}

func createConfirmedUser(db *gorm.DB) *User {
	now := time.Now().Add(2 * time.Minute)
	u := &User{
		Email: fmt.Sprintf("sdsuhb%v@asdji.co", now.UnixNano()),
		Password: "adsjfasdfasdfa",
		ConfirmedAt: &now,
	}
	err := u.Create(db)
	if err != nil {
		panic(err)
	}
	return u
}
