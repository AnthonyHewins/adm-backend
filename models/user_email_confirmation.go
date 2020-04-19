package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Must be within 15 minutes
const confirmationTooLate = 15

type UserEmailConfirmation struct {
	UserId int64        `gorm:"primary_key; auto_increment: false"`
	Token  string       `gorm:"primary_key; auto_increment: false"`
}

func (uec *UserEmailConfirmation) ConfirmEmail(db *gorm.DB) error {
	if uec.UserId <= 0 {
		if err := db.Debug().Where("token = ?", uec.Token).First(&uec).Error; err != nil {
			return err
		}
	}

	u := User{ID: uec.UserId}
	if err := db.First(&u).Error; err != nil { return err }

	if time.Now().Sub(u.RegisteredAt).Minutes() > confirmationTooLate {
		return &EmailConfirmationLate
	}

	// DB trigger deletes (all) the UserEmailConfirmation(s) for this user upon this update happening
	err := db.Exec("update users set confirmed_at = CURRENT_TIMESTAMP where id = ?", u.ID).Error
	if err != nil { return err }

	return nil
}
