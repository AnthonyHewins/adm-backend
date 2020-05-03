package models

import (
	"time"

	"github.com/AnthonyHewins/adm-backend/smtp"
	"github.com/jinzhu/gorm"
)

// Hacky way of upserting due to gorm's lack of upsert at the moment
const upsertQuery = "INSERT INTO user_password_resets (user_id, token, reset_at) VALUES (?, ?, CURRENT_TIMESTAMP) ON CONFLICT (user_id) DO UPDATE SET token = ?, reset_at = CURRENT_TIMESTAMP"

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

	return db.Transaction(func(tx *gorm.DB) error {
		tx.Exec(upsertQuery, upr.UserID, token, token)

		if err != nil { return err }

		return smtp.PasswordReset(u.Email, token)
	})
}
