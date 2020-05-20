package models

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DB struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SslMode  string `yaml:"sslmode"`
}

var (
	masterConfig string

	TokenTimeout = &Error{s: "token has expired; you will need a new one to proceed with this action"}
)

type Error struct{ s string }

func (e *Error) Error() string {
	return e.s
}

func Connect() (*gorm.DB, error) {
	return gorm.Open("postgres", masterConfig)
}

func DBSetup(dbConfig *DB) {
	var sslmode string
	if (dbConfig.SslMode == "") {
		sslmode = "require"
	} else {
		sslmode = dbConfig.SslMode
	}

	masterConfig = fmt.Sprintf(
		"host=%v port=%v dbname=%v user=%v password=%v sslmode=%v",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		dbConfig.User,
		dbConfig.Password,
		sslmode,
	)
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
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b)[:40], nil
}
