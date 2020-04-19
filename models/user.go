package models

import (
	"fmt"
	"time"
	"math/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"

	"github.com/AnthonyHewins/adm-backend/smtp"
)

const passwordCost = bcrypt.DefaultCost + 2

type User struct {
	ID           int64
	Email        string     `gorm:"type:varchar(255); not null; unique"`
	Password     string     `gorm:"type:varchar(60);  not null"`
	RegisteredAt time.Time  `gorm:"not null"`
	ConfirmedAt  *time.Time
}

func CreateUser(db *gorm.DB, email, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	if err != nil { return nil, err }

	// do this immediately; plaintext should be in memory as short as possible
	// this will help push the go runtime to drop pw
	password = string(hash)
	u := User{ Email: email, Password: password }

	token, err := base64ConfirmationString()
	if err != nil { return nil, err }

	return &u, db.Transaction(func (t *gorm.DB) error {
		err := t.Create(&u).Error
		if err != nil { return err }

		err = t.Create(&UserEmailConfirmation{
			UserId: u.ID,
			Token: token,
		}).Error

		if err != nil { return err }

		return smtp.AccountConfirmation(u.Email, token)
	})
}

func (u *User) RefreshConfirmationToken(db *gorm.DB) error {
	if u.ConfirmedAt != nil { return fmt.Errorf("already confirmed email") }

	return db.Transaction(func (tx *gorm.DB) error {
		err := tx.Exec("update users set registered_at = CURRENT_TIMESTAMP where id = ?", u.ID).Error
		if err != nil { return err }

		token, err := base64ConfirmationString()
		if err != nil { return err }

		err = tx.Exec("update user_email_confirmations set token = ? where user_id = ?", token, u.ID).Error

		if err != nil { return err }

		return smtp.AccountConfirmation(u.Email, token)
	})
}

func base64ConfirmationString() (string, error) {
	// More random in the future?
	// 1. csprng
	// 2. Add PID
	// 3. Salt with username/email
	rand.Seed(time.Now().UnixNano())

	// Atrociously bad algo ATM, but good for now;
	// improve the speed and wasted computations later
	b := make([]byte, 40, 40)
	if _, err := rand.Read(b); err != nil { return "", err }

	return base64.URLEncoding.EncodeToString(b)[:40], nil
}
