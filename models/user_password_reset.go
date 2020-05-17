package models

/*
import (
	"time"

	"github.com/AnthonyHewins/adm-backend/smtp"
	"github.com/jinzhu/gorm"
)

// Hacky way of upserting due to gorm's lack of upsert at the moment
const upsertQuery = `
INSERT INTO user_password_resets (user_id, token, reset_at)
VALUES (?, ?, CURRENT_TIMESTAMP) ON CONFLICT (user_id)
DO UPDATE SET token = ?, reset_at = CURRENT_TIMESTAMP
`

type UserPasswordReset struct {
	UserID     uint64
	Token      string
	ResetAt    time.Time
}

func (upr *UserPasswordReset) CreateResetPasswordToken(db *gorm.DB) error {
	token, err := base64ConfirmationString()
	if err != nil { return err }

	u := User{ID: upr.UserID}
	if err := db.First(&u).Error; err != nil { return err }

	upr.Token = token

	return db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Exec(upsertQuery, upr.UserID, token, token).Error; err != nil {
			return err
		}

		return smtp.PasswordReset(u.Email, token)
	})
}

func (upr *UserPasswordReset) ResetPassword(db *gorm.DB, newPw string) error {
	var err error
	if upr.UserID == 0 {
		err = db.Where("token = ?", upr.Token).First(&upr).Error
	} else {
		err = db.Where("user_id = ?", upr.UserID).First(&upr).Error
	}

	if err != nil { return err }

	if userTokenExpired(upr.ResetAt) {
		return TokenTimeout
	}

	u := User{ID: upr.UserID, Password: newPw}
	return u.ResetPassword(db)
}
*/
