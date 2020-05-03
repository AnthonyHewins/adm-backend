package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
)

func TestResetPassword(t *testing.T) {
	DBSetupTest(nil)
	db, _ := Connect()

	// Nonexiting user -> error
	upr := UserPasswordReset{UserID: 10000000000}
	err := upr.CreateResetPasswordToken(db)
	if !gorm.IsRecordNotFoundError(err) {
		t.Errorf("should not have found the user record, but something else happened: %v", err)
	}

	// Existing user requesting PW reset for the first time -> create new PW reset entry in DB
	u, _ := CreateUser(db, fmt.Sprintf("asod%v@soijdf.co", time.Now().UnixNano()), "sdjnsdfn")
	upr.UserID = u.ID
	upr.CreateResetPasswordToken(db)
	assertUPR(db, t, upr)

	// User resets the timer on the PW reset -> update the PW reset entry in the DB
	db.Where("user_id = ?", upr.UserID).First(&upr)
	upr.CreateResetPasswordToken(db)
	assertUPR(db, t, upr)
}

func assertUPR(db *gorm.DB, t *testing.T, mustBeDifferentThan UserPasswordReset) {
	upr := UserPasswordReset{}
	err := db.Where("user_id = ?", mustBeDifferentThan.UserID).First(&upr).Error
	if err != nil {
		t.Errorf(
			"UserPasswordReset should've existed with the same ID as this struct (%v) but all other fields should have been changed. Instead got error: %v",
			mustBeDifferentThan,
			err,
		)
	}

	// If the caller passes anything but "" to this function they want to confirm that
	// the database should have reflected a change in the token for this particular user
	if mustBeDifferentThan.Token != "" {
		if upr.Token == mustBeDifferentThan.Token {
			t.Errorf("Token should have been updated")
		}
	}

	// If the caller passes anything but time.Zero to this function they want to confirm that
	// the database should have reflected a change in the updated time for this particular user
	if !mustBeDifferentThan.ResetAt.IsZero() {
		if upr.ResetAt == mustBeDifferentThan.ResetAt {
			t.Errorf("Token should have been updated")
		}
	}
}
