package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Must be within 15 minutes
const confirmationThreshold = 15 * time.Minute

type UserEmailConfirmation struct {
	UserID uint64
	Token  string
}

func (uec *UserEmailConfirmation) ConfirmEmail(db *gorm.DB) error {
	if uec.UserID <= 0 {
		if err := db.Where("token = ?", uec.Token).First(&uec).Error; err != nil {
			return err
		}
	}

	u := User{ID: uec.UserID}
	if err := db.First(&u).Error; err != nil { return err }

	minutesPassedSinceRegistration := time.Duration(time.Since(u.RegisteredAt).Minutes())

	if minutesPassedSinceRegistration > confirmationThreshold {
		return &EmailConfirmationLate
	}

	// DB trigger deletes (all) the UserEmailConfirmation(s) for this user upon this update happening
	err := db.Exec("update users set confirmed_at = CURRENT_TIMESTAMP where id = ?", u.ID).Error
	if err != nil { return err }

	return nil
}
