package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"

	"github.com/AnthonyHewins/adm-backend/smtp"
)

const (
	passwordCost = bcrypt.DefaultCost + 2

	passwordLength   = 6
	emailRegexString = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

var (
	emailRegex        = regexp.MustCompile(emailRegexString)
	InvalidEmail      = &Error{s: "email didn't pass validation"}
	PasswordTooSimple = &Error{s: fmt.Sprintf("password needs to be at least %v characters", passwordLength)}
	AlreadyConfirmed  = &Error{s: "email is already confirmed"}
	EmailNotConfirmed = &Error{s: "you have not confirmed your email yet; please confirm it to log in"}
)

type User struct {
	ID             uint64
	Email          string    `gorm:"type: varchar(255); not null; unique"`
	Password       string    `gorm:"-"`
	HashedPassword string    `gorm:"column:password; type: varchar(60); not null"`
	RegisteredAt   time.Time `gorm:"default: CURRENT_TIMESTAMP; not null"`
	ConfirmedAt    *time.Time
}

func (u *User) Create(db *gorm.DB) error {
	if !isPasswordValid(u.Password) {
		return PasswordTooSimple
	}

	if !emailRegex.MatchString(u.Email) {
		return InvalidEmail
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), passwordCost)
	if err != nil {
		return err
	}

	token, err := base64ConfirmationString()
	if err != nil {
		return err
	}

	u.Password = ""
	u.HashedPassword = string(hash)
	u.Email = strings.ToLower(u.Email)

	return db.Transaction(func(t *gorm.DB) error {
		err := t.Create(&u).Error
		if err != nil {
			return err
		}

		err = t.Create(&UserEmailConfirmation{UserID: u.ID, Token: token}).Error
		if err != nil {
			return err
		}

		return smtp.AccountConfirmation(u.Email, token)
	})
}

func (u *User) RefreshConfirmationToken(db *gorm.DB) (string, error) {
	if u.ConfirmedAt != nil {
		return "", AlreadyConfirmed
	}

	token, err := base64ConfirmationString()
	if err != nil {
		return token, err
	}

	return token, db.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec("update users set registered_at = CURRENT_TIMESTAMP where id = ?", u.ID).Error
		if err != nil {
			return err
		}

		return tx.Exec("update user_email_confirmations set token = ? where user_id = ?", token, u.ID).Error
	})
}

func (u *User) Authenticate(db *gorm.DB) error {
	var err error
	if u.ID > 0 {
		err = db.First(&u).Error
	} else {
		err = db.Where("email = ?", u.Email).First(&u).Error
	}

	if err != nil {
		return err
	}

	// Important: check the auth first. If they aren't confirmed and can't log
	// in anyways, we need to not give any extra information
	err = bcrypt.CompareHashAndPassword(
		[]byte(u.HashedPassword),
		[]byte(u.Password),
	)

	if err != nil {
		return err
	}

	if u.ConfirmedAt == nil {
		return EmailNotConfirmed
	}

	return nil
}

func (u *User) ResetPassword(db *gorm.DB) error {
	if !isPasswordValid(u.Password) {
		return PasswordTooSimple
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), passwordCost)
	if err != nil {
		return err
	}

	u.HashedPassword = string(hash)
	u.Password = ""

	return db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Model(&u).Update("password", u.HashedPassword).Error; err != nil {
			return err
		}

		return tx.Where("user_id = ?", u.ID).Delete(&UserPasswordReset{}).Error
	})
}

func isPasswordValid(s string) bool {
	// this will grow in the future
	return len(s) >= passwordLength
}

func userTokenExpiryCheck(t time.Time) bool {
	return time.Since(t).Minutes() < TokenTimeoutThreshold
}
