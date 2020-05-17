package models

/*
import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCreateResetPasswordToken(t *testing.T) {
	DBSetupTest(nil)
	db, _ := Connect()

	// Nonexiting user -> error
	upr := UserPasswordReset{UserID: 10000000000}
	err := upr.CreateResetPasswordToken(db)
	if !gorm.IsRecordNotFoundError(err) {
		t.Errorf("should not have found the user record, but something else happened: %v", err)
	}

	// Existing user requesting PW reset for the first time -> create new PW reset entry in DB
	u := createUser(db)
	upr.UserID = u.ID
	upr.CreateResetPasswordToken(db)
	assertUPR(db, t, UserPasswordReset{UserID: u.ID})

	// User resets the timer on the PW reset -> update the PW reset entry in the DB
	copyOfA := upr

	db.Where("user_id = ?", upr.UserID).First(&upr)
	upr.CreateResetPasswordToken(db)
	assertUPR(db, t, copyOfA)
}

func TestResetPassword(t *testing.T) {
	db := getDB(t)

	upr := UserPasswordReset{}

	// UserID, Token not found -> RecordNotFound
	upr.UserID = 99999999
	assert.Equal(t, gorm.ErrRecordNotFound, upr.ResetPassword(db, ""))
	upr.Token = "ddd"
	assert.Equal(t, gorm.ErrRecordNotFound, upr.ResetPassword(db, ""))

	u := createConfirmedUser(db)
	upr.UserID = u.ID
	if err := upr.CreateResetPasswordToken(db); err != nil {
		t.Fatal(err)
	}

	oldTokenTimeoutValue := TokenTimeoutThreshold

	// Token expired -> TokenTimeout
	TokenTimeoutThreshold = 0
	assert.Equal(t, TokenTimeout, upr.ResetPassword(db, "fhusdh"))

	TokenTimeoutThreshold = oldTokenTimeoutValue

	// Weak pw -> PasswordTooSimple
	assert.Equal(t, PasswordTooSimple, upr.ResetPassword(db, "sdd"))

	// Proper PW -> reset is gone through
	assert.Equal(t, nil, upr.ResetPassword(db, "aujsdnuajsdnjuajs"))
}

func assertUPR(db *gorm.DB, t *testing.T, mustBeDifferentThan UserPasswordReset) {
	upr := UserPasswordReset{}
	err := db.Where("user_id = ?", mustBeDifferentThan.UserID).First(&upr).Error
	assert.Equal(t, nil, err)

	// If the caller passes anything but "" to this function they want to confirm that
	// the database should have reflected a change in the token for this particular user
	if mustBeDifferentThan.Token != "" {
		assert.NotEqual(t, mustBeDifferentThan.Token, upr.Token)
	}

	// If the caller passes anything but time.Zero to this function they want to confirm that
	// the database should have reflected a change in the updated time for this particular user
	if !mustBeDifferentThan.ResetAt.IsZero() {
		assert.NotEqual(t, mustBeDifferentThan.ResetAt, upr.ResetAt)
	}
}
*/
