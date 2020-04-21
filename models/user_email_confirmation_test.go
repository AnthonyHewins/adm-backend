package models

import (
	"fmt"
	"time"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestEmailConfirmation(t *testing.T) {
	DBSetupTest(nil)

	db, err := Connect()
	if err != nil { t.Fatal(err.Error()) }

	u, err := CreateUser(db, fmt.Sprintf("a%v@asdijf.com", time.Now().UnixNano()), "sadasf")
	if err != nil { t.Fatal(err.Error()) }

	// Test instances where the confirmation isn't found
	notfound(
		t,
		db,
		UserEmailConfirmation{}, // blank
		UserEmailConfirmation{Token: "nonsense"},
		UserEmailConfirmation{UserID: 100000},
	)

	// Fetch uec from the user that was inserted at beginning of tests
	// (inserted from CreateUser)
	uec := UserEmailConfirmation{}
	err = db.Where("user_id = ?", u.ID).First(&uec).Error
	if err != nil { t.Fatalf(err.Error()) }

	// Test based on using UserID to immediately find the user
	userid(t, db, u, &uec)

	// Test based on finding the user from the Token
	token(t, db, u, &uec)
}



func notfound(t *testing.T, db *gorm.DB, uecs ...UserEmailConfirmation) {
	for _, uec := range uecs {
		err := uec.ConfirmEmail(db)
		if err == nil || err.Error() != "record not found" {
			t.Errorf("should have not found empty record and got RecordNotFound, got %v", err)
		}
	}
}



func token(t *testing.T, db *gorm.DB, u *User, actualUec *UserEmailConfirmation) {
	uec := UserEmailConfirmation{Token: actualUec.Token}

	// UserId = 0 -> fetches UserId -> checks against, works
	if err := uec.ConfirmEmail(db); err != nil {
		t.Errorf("should have successfully confirmed, got %v", err)
	}

	// Check some assertions: confirmed_at should not be nil, UserEmailConfirmation
	// should have been deleted now.
	// Reinsert the rows that were deleted due to the email being confirmed
	checkAndReinsert(t, db, u, actualUec)
}



func userid(t *testing.T, db *gorm.DB, u *User, actualUec *UserEmailConfirmation) {
	uec := UserEmailConfirmation{UserID: actualUec.UserID}

	// UserID exists -> immediately tries confirm
	if err := uec.ConfirmEmail(db); err != nil {
		t.Errorf("should have successfully confirmed, got %v", err)
	}

	// Check some assertions: confirmed_at should not be nil, UserEmailConfirmation
	// should have been deleted now.
	// Reinsert the rows that were deleted due to the email being confirmed
	checkAndReinsert(t, db, u, actualUec)
}



func checkAndReinsert(t *testing.T, db *gorm.DB, u *User, uec *UserEmailConfirmation) {
	userInDb := User{ID: u.ID}
	db.First(&userInDb)

	if userInDb.ConfirmedAt == nil {
		t.Errorf("user should be confirmed but was not, got %v for ConfirmedAt", userInDb.ConfirmedAt)
	}

	anyUec := UserEmailConfirmation{}
	if query := db.Where("user_id = ?", u.ID).First(&anyUec); !query.RecordNotFound() {
		t.Errorf("user %v should have no tokens left but does: %v", userInDb, anyUec)
	}


	err := db.Model(u).Where("id = ?", u.ID).Update("confirmed_at", nil).Error

	if err != nil {
		t.Fatalf(err.Error())
	}

	// Recreate the UEC exactly as it was, same token and ID
	err = db.Create(uec).Error

	if err != nil {
		t.Fatalf(err.Error())
	}
}
