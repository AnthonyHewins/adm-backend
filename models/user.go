package models

import (
	"fmt"
	"time"
	"regexp"
	"math/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
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

	// do this immediately; plaintext should be in memory as short as possible
	// this will help push the go runtime to drop pw
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

		return t.Create(&UserEmailConfirmation{ UserID: u.ID, Token: token }).Error
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
