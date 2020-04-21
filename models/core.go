package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DB struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

var (
	masterConfig string

	EmailConfirmationLate = Error{s: "email confirmation expired"}
)

type Error struct { s string }

func (e *Error) Error() string {
	return e.s
}

func Connect() (*gorm.DB, error) {
	return gorm.Open("postgres", masterConfig)
}

func DBSetup(dbConfig *DB) {
	masterConfig = fmt.Sprintf(
		"host=%v port=%v dbname=%v user=%v password=%v",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		dbConfig.User,
		dbConfig.Password,
	)
}
