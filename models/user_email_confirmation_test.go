package models

import (
	"testing"

	"github.com/jinzhu/gorm"
)

func TestUserEmailConfirmation_ConfirmEmail(t *testing.T) {
	db, _ := Connect()
	defer db.Close()

	willConfirm := createUser(db)

	tests := []struct {
		name    string
		uec UserEmailConfirmation
		err error
	}{
		{
			"Not found due to token being meaningless",
			UserEmailConfirmation{Token: "3u49823"},
			gorm.ErrRecordNotFound,
		},
		{
			"Not found due to UserID being meaningless",
			UserEmailConfirmation{UserID: 490283482039},
			gorm.ErrRecordNotFound,
		},
		{
			"User is confirmed",
			UserEmailConfirmation{UserID: willConfirm.ID},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.uec.ConfirmEmail(db); err != tt.err {
				t.Errorf("UserEmailConfirmation.ConfirmEmail() error = %v, wantErr %v", err, tt.err)
			}
		})
	}

	uec := UserEmailConfirmation{}
	err := db.Where("user_id = ?", willConfirm.ID).First(&uec)
	if !err.RecordNotFound() {
		t.Errorf("all UECs should have been deleted and the weren't, found one: %v", err)
	}
}
