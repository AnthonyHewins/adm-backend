package models

import (
	"github.com/jinzhu/gorm"
)

// Must be within 15 minutes
var TokenTimeoutThreshold = float64(15)

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

	if userTokenExpiryCheck(u.RegisteredAt) {
		return TokenTimeout
	}

	// DB trigger deletes (all) the UserEmailConfirmation(s) for this user upon this update happening
	err := db.Exec("update users set confirmed_at = CURRENT_TIMESTAMP where id = ?", u.ID).Error
	if err != nil { return err }

	return nil
}
