package controllers

import "github.com/AnthonyHewins/adm-backend/models"

type credentials struct {
	Email string    `json:"email"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *credentials) toUser() *models.User {
	return &models.User{
		Email: c.Email,
		Password: c.Password,
	}
}
