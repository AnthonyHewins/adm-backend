package models

import (
	"fmt"
	"time"
	"regexp"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"

	"github.com/AnthonyHewins/adm-backend/smtp"
)

const (
	passwordCost     = bcrypt.DefaultCost + 2

	passwordLength   = 6
	emailRegexString = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

var (
	emailRegex        = regexp.MustCompile(emailRegexString)
	InvalidEmail      = &Error{s: "email didn't pass validation"}
	PasswordTooSimple = &Error{s: fmt.Sprintf("password needs to be at least %v characters", passwordLength)}
	AlreadyConfirmed  = &Error{s: "email is already confirmed"}
)

type User struct {
	ID           uint64
	Email        string     `gorm:"type:         varchar(255); not null; unique"`
	Password     string     `gorm:"type:          varchar(60); not null"`
	RegisteredAt time.Time  `gorm:"default: CURRENT_TIMESTAMP; not null"`
	ConfirmedAt  *time.Time
}

func CreateUser(db *gorm.DB, email, password string) (*User, error) {
	if !isPasswordValid(password) {
		return nil, PasswordTooSimple
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	if err != nil { return nil, err }

	password = string(hash)
	u := User{ Email: email, Password: password }

	if !emailRegex.MatchString(email) {
		return &u, InvalidEmail
	}

	token, err := base64ConfirmationString()
	if err != nil { return nil, err }

	return &u, db.Transaction(func (t *gorm.DB) error {
		err := t.Create(&u).Error
		if err != nil { return err }

		err = t.Create(&UserEmailConfirmation{ UserID: u.ID, Token: token }).Error
		if err != nil { return err }

		return smtp.AccountConfirmation(email, token)
	})
}

func (u *User) RefreshConfirmationToken(db *gorm.DB) (string, error) {
	if u.ConfirmedAt != nil { return "", AlreadyConfirmed }

	token, err := base64ConfirmationString()
	if err != nil { return token, err }

	return token, db.Transaction(func (tx *gorm.DB) error {
		err := tx.Exec("update users set registered_at = CURRENT_TIMESTAMP where id = ?", u.ID).Error
		if err != nil { return err }

		return tx.Exec("update user_email_confirmations set token = ? where user_id = ?", token, u.ID).Error
	})
}

func isPasswordValid(s string) bool {
	// this will grow in the future
	return len(s) >= passwordLength
}
