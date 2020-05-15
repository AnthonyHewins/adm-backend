package auth

import (
	"github.com/AnthonyHewins/adm-backend/models"
	"github.com/jinzhu/gorm"
)

func dbInstance() *gorm.DB {
	models.DBSetup(&models.DB{
		Host: "localhost",
		Port: 5432,
		Name: "admtest",
		User: "test",
		Password: "test",
	})

	db, err := models.Connect()

	if err != nil {
		panic(err)
	}

	return db
}
